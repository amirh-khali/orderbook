package kafka

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"orderbook/db/models"
)

var consumer *kafka.Consumer

func InitConsumer(servers string, groupID string, reset string) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": servers,
		"group.id":          groupID,
		"auto.offset.reset": reset,
	})

	if err != nil {
		log.Printf("failed to create consumer: %s\n", err)
		os.Exit(1)
	}

	consumer = c
	log.Printf("consumer created %v\n", consumer)
}

func Subscribe(t string) {
	err := consumer.Subscribe(t, nil)
	if err != nil {
		log.Printf("failed to subscribe %s topic\n", t)
		os.Exit(1)
	}
	log.Printf("topic %s subscribed\n", t)
}

func StartConsume() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	run := true
	for run {
		select {
		case sig := <-ch:
			fmt.Printf("Caught signal %v: terminating\n", sig)
			run = false
		default:
			ev := consumer.Poll(100)
			if ev == nil {
				continue
			}

			switch e := ev.(type) {
			case *kafka.Message:
				var value models.Order
				err := json.Unmarshal(e.Value, &value)
				if err != nil {
					fmt.Printf("Failed to deserialize payload: %s\n", err)
				} else {
					fmt.Printf("%% Message on %s:\n%+v\n", e.TopicPartition, value)
				}
				if e.Headers != nil {
					fmt.Printf("%% Headers: %v\n", e.Headers)
				}
			case kafka.Error:
				_, _ = fmt.Fprintf(os.Stderr, "%% Error: %v: %v\n", e.Code(), e)
			default:
				fmt.Printf("Ignored %v\n", e)
			}
		}
	}

	fmt.Printf("Closing consumer\n")
	_ = consumer.Close()
}
