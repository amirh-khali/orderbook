package kafka

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"orderbook/core"
	kafkamodels "orderbook/kafka/models"
)

var consumer *kafka.Consumer

func InitConsumer(servers string, groupID string, reset string) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": servers,
		"group.id":          groupID,
		"auto.offset.reset": reset,
	})

	if err != nil {
		log.Fatalf("failed to create consumer: %s\n", err)
	}

	consumer = c
	log.Printf("consumer created %v\n", consumer)
}

func Subscribe(t string) {
	err := consumer.Subscribe(t, nil)
	if err != nil {
		log.Fatalf("failed to subscribe %s topic\n", t)
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
			log.Printf("caught signal %v: terminating\n", sig)
			run = false
		default:
			ev := consumer.Poll(100)
			if ev == nil {
				continue
			}

			switch e := ev.(type) {
			case *kafka.Message:
				var or kafkamodels.OrderRequest
				err := json.Unmarshal(e.Value, &or)
				if err != nil {
					log.Printf("failed to deserialize payload: %s\n", err)
				} else {
					log.Printf("%% message on %s: %+v\n", e.TopicPartition, or)
					(*core.OrderbookMap)[or.Symbol].AddOrder(kafkamodels.NewOrder(or))
				}
				if e.Headers != nil {
					log.Printf("%% headers: %v\n", e.Headers)
				}
			case kafka.Error:
				log.Printf("%% error: %v: %v\n", e.Code(), e)
			default:
				log.Printf("ignored %v\n", e)
			}
		}
	}

	log.Printf("closing consumer\n")
	_ = consumer.Close()
}
