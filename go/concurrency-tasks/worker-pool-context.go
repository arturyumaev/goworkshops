package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

func randomTimeWork(id string) error {
	dur := time.Duration(rand.Intn(1000))
	time.Sleep(dur * time.Millisecond)

	if dur < 100 {
		return errors.New("error")
	}

	return nil
}

func workerPoolContext() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	numCPU := runtime.NumCPU()
	wg := &sync.WaitGroup{}
	wg.Add(numCPU)

	work := make(chan string)
	errs := make(chan error)

	go func() {
		wg.Wait()
		close(errs)
	}()

	go func() {
		defer close(work)

		for t := 0; t < 1000; t++ {
			select {
			case <-ctx.Done():
				return
			default:
				work <- fmt.Sprint(t)
			}
		}
	}()

	for i := 0; i < numCPU; i++ {
		go func() {
			defer wg.Done()

			for url := range work {
				err := randomTimeWork(url)

				select {
				case <-ctx.Done():
					return
				default:
					errs <- err
				}
			}
		}()
	}

	for err := range errs {
		if err != nil {
			cancel()
		}

		fmt.Println(err)
	}
}
