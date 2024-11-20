package main

import (
	"encoding/json"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/gadisamenu/tolling/types"
)

type DataProducer interface {
	ProduceData(types.ObuData) error
}

type KafkaProducer struct {
	topic    string
	producer *kafka.Producer
}

func NewKafkaProducer(topic string) (*KafkaProducer, error) {
	prod, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost"})
	if err != nil {
		return nil, err
	}

	// // Delivery report handler for produced messages
	// go func() {
	// 	for e := range prod.Events() {
	// 		switch ev := e.(type) {
	// 		case *kafka.Message:
	// 			if ev.TopicPartition.Error != nil {
	// 				// fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
	// 			} else {
	// 				// fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
	// 			}
	// 		}
	// 	}
	// }()

	return &KafkaProducer{
		topic:    topic,
		producer: prod,
	}, nil
}

func (pr *KafkaProducer) ProduceData(data types.ObuData) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	pr.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &pr.topic, Partition: kafka.PartitionAny},
		Value:          b,
	}, nil)
	return nil
}
