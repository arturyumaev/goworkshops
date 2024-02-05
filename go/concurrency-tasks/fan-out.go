package main

import (
	"fmt"
	"sync"
)

func fanOut() {
	errs := make(chan error)

	wg := &sync.WaitGroup{}
	wg.Add(len(urls))

	go func() {
		wg.Wait()
		close(errs)
	}()

	for _, url := range urls {
		url := url

		go func() {
			defer wg.Done()

			err := healthCheck(url)
			errs <- err
		}()
	}

loop:
	for {
		select {
		case err, open := <-errs:
			if !open {
				break loop
			}

			fmt.Println(err)
		}
	}
}
