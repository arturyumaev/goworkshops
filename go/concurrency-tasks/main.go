package main

import (
	"fmt"
	"math/rand"
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
	time.Sleep(time.Duration(rand.Intn(2000)) * time.Millisecond)
	return nil
}

func main() {
	fmt.Println("fan-out")
	fanOut()

	fmt.Println("limited fan-out")
	limitedFanOut()

	fmt.Println("worker pool")
	workerPool()

	fmt.Println("worker pool context")
	workerPoolContext()
}
