package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

var (
	urls = []string{
		"http://node-1:3000",
		"http://node-2:3000",
		"http://node-3:3000",
		"http://node-4:3000",
		"http://node-5:3000",
	}
)

func healthCheck(url string) error {
	url += "/healthcheck"
	time.Sleep(time.Duration(rand.Intn(2000)) * time.Millisecond)
	return nil
}

/*
// Example

func healthCheck(url string) error {
	url += "/healthcheck"
	time.Sleep(time.Second)
	return nil
}

func main() {
	for _, url := range urls {
		err := healthCheck(url)
		fmt.Println(err)
	}
}
*/

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

	for err := range errs {
		fmt.Println(err)
	}
}

func limitedFanOut() {
	simGoroutines := 3
	blockChan := make(chan struct{}, simGoroutines)
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

func workerPool() {
	N := runtime.NumCPU()
	in := make(chan string)
	errs := make(chan error, len(urls))

	wg := &sync.WaitGroup{}
	wg.Add(len(urls))

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

	for i := 0; i < N; i++ {
		go func() {
			for url := range in {
				err := healthCheck(url)
				errs <- err
				wg.Done()
			}
		}()
	}

	for err := range errs {
		fmt.Println(err)
	}
}

func main() {
	fmt.Println("fan-out")
	fanOut()
	fmt.Println("limited fan-out")
	limitedFanOut()
	fmt.Println("worker pool")
	workerPool()
}
