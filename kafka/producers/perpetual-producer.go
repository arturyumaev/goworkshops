package main

import (
	"encoding/json"
	"errors"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/IBM/sarama"
)

type entity struct {
	ID   []byte
	Name string
}

func (e *entity) Bytes() []byte {
	bytes, _ := json.Marshal(&e)
	return bytes
}

func newEntity() *entity {
	return &entity{
		ID:   []byte(strconv.Itoa(rand.Intn(100))),
		Name: "Test",
	}
}

var (
	logger  = log.Default()
	brokers = []string{
		"localhost:19092", // broker_id=1
		"localhost:29092", // broker_id=2
		"localhost:39092", // broker_id=3
	}
	topic = "test"
)

func createTestTopic(cfg *sarama.Config) {
	admin, err := sarama.NewClusterAdmin(brokers, cfg)
	if err != nil {
		logger.Fatalln(err)
	}

	defer func() {
		_ = admin.Close()
	}()

	err = admin.CreateTopic(
		topic,
		&sarama.TopicDetail{
			NumPartitions:     3,
			ReplicationFactor: 3,
		},
		false,
	)
	if err != nil && errors.Unwrap(err) != sarama.ErrTopicAlreadyExists {
		logger.Fatalln(err)
	}
}

func deleteTestTopic(cfg *sarama.Config) {
	admin, err := sarama.NewClusterAdmin(brokers, cfg)
	if err != nil {
		logger.Fatalln(err)
	}

	defer func() {
		_ = admin.Close()
	}()

	if err := admin.DeleteTopic(topic); err != nil {
		logger.Fatalln(err)
	}
}

func newProducer(cfg *sarama.Config) sarama.SyncProducer {
	cfg.Producer.Partitioner = sarama.NewRandomPartitioner
	cfg.Producer.RequiredAcks = sarama.WaitForLocal
	cfg.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokers, cfg)
	if err != nil {
		logger.Fatalln(err)
	}

	return producer
}

func main() {
	config := sarama.NewConfig()
	createTestTopic(config)
	producer := newProducer(config)

	go func() {
		done := make(chan os.Signal, 1)
		signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
		<-done

		producer.Close()
		deleteTestTopic(config)
		log.Println("graceful shutdown...")
		os.Exit(1)
	}()

	for {
		e := newEntity()

		msg := &sarama.ProducerMessage{
			Topic: topic,
			Key:   sarama.ByteEncoder(e.ID),
			Value: sarama.ByteEncoder(e.Bytes()),
		}

		_, _, err := producer.SendMessage(msg)
		if err == nil {
			logger.Printf("sent message '%s' to topic '%s'", e.ID, topic)
		}

		randomMsSleep := rand.Intn(1000-100) + 100 // generates random int in range [100, 1000)
		time.Sleep(time.Duration(randomMsSleep) * time.Millisecond)
	}
}
