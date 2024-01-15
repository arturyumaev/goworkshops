package main

import (
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/IBM/sarama"
)

var (
	maxMessages = 6_000_000
	successes   int32
	errors      int32
)

func main() {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForLocal
	config.Producer.Return.Successes = true

	producer, err := sarama.NewAsyncProducer([]string{"localhost:19092"}, config)
	if err != nil {
		panic(err)
	}

	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for range producer.Successes() {
			atomic.AddInt32(&successes, 1)
		}
	}()

	go func() {
		defer wg.Done()
		for err := range producer.Errors() {
			atomic.AddInt32(&errors, 1)
			fmt.Println(err)
		}
	}()

	for i := 1; i <= maxMessages; i++ {
		message := sarama.ProducerMessage{Topic: "test", Value: sarama.StringEncoder("testing 123")}
		producer.Input() <- &message

		if i > 0 && i%100 == 0 {
			fmt.Printf("wrote %d messages\n", i)
		}
	}
	go producer.AsyncClose()

	wg.Wait()

	fmt.Printf("Stats: successes: %d, errors: %d\n", successes, errors)
}
