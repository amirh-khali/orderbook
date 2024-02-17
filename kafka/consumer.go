package kafka

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"orderbook/core"
	"orderbook/core/models"
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
			log.Printf("Caught signal %v: terminating\n", sig)
			run = false
		default:
			ev := consumer.Poll(100)
			if ev == nil {
				continue
			}

			switch e := ev.(type) {
			case *kafka.Message:
				var o models.Order
				err := json.Unmarshal(e.Value, &o)
				if err != nil {
					log.Printf("Failed to deserialize payload: %s\n", err)
				} else {
					log.Printf("%% Message on %s: %+v\n", e.TopicPartition, o)
					o.Renew()
					(*core.OrderbookMap)[o.Symbol].AddOrder(&o)
				}
				if e.Headers != nil {
					log.Printf("%% Headers: %v\n", e.Headers)
				}
			case kafka.Error:
				log.Printf("%% Error: %v: %v\n", e.Code(), e)
			default:
				log.Printf("Ignored %v\n", e)
			}
		}
	}

	log.Printf("Closing consumer\n")
	_ = consumer.Close()
}
