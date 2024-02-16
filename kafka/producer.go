package kafka

import (
	"encoding/json"
	"log"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

var producer *kafka.Producer

func InitProducer(servers string) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": servers,
	})

	if err != nil {
		log.Printf("failed to create producer: %s\n", err)
		os.Exit(1)
	}

	producer = p
	log.Printf("producer created %v\n", producer)
}

func Produce(msg any, topic string) {
	deliveryChan := make(chan kafka.Event)

	payload, err := json.Marshal(msg)
	if err != nil {
		log.Fatal(err)
	}

	err = producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          payload,
	}, deliveryChan)

	if err != nil {
		log.Printf("produce failed: %v\n", err)
		os.Exit(1)
	}

	e := <-deliveryChan
	m := e.(*kafka.Message)

	if m.TopicPartition.Error != nil {
		log.Printf("delivery failed: %v\n", m.TopicPartition.Error)
	} else {
		log.Printf("delivered message to topic %s [%d] at offset %v\n",
			*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
	}

	close(deliveryChan)
}
