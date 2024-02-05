package main

import (
	"fmt"
	"runtime"
	"sync"
)

func limitedFanOut() {
	blockChan := make(chan struct{}, runtime.NumCPU())
	defer close(blockChan)
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
			blockChan <- struct{}{}
			defer wg.Done()

			err := healthCheck(url)
			errs <- err

			<-blockChan
		}()
	}

	for err := range errs {
		fmt.Println(err)
	}
}
