package main

import (
	"fmt"
	"runtime"
	"sync"
)

func workerPool() {
	in := make(chan string)
	errs := make(chan error)

	workers := runtime.NumCPU()
	wg := &sync.WaitGroup{}
	wg.Add(workers)

	go func() {
		wg.Wait()
		close(errs)
	}()

	go func() {
		for _, url := range urls {
			in <- url
		}
		close(in)
	}()

	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()

			for url := range in {
				err := healthCheck(url)
				errs <- err
			}
		}()
	}

	for err := range errs {
		fmt.Println(err)
	}
}
