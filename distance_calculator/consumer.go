package main

import (
	"encoding/json"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/gadisamenu/tolling/types"
	"github.com/sirupsen/logrus"
)

type KafkaConsumer struct {
	consumer    *kafka.Consumer
	isRunning   bool
	calcService CalculatorServicer
}

func NewKafkaConsumer(topic string, service CalculatorServicer) (*KafkaConsumer, error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost",
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		return nil, err
	}

	err = c.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		return nil, err
	}

	return &KafkaConsumer{
		consumer:    c,
		calcService: service,
	}, nil
}

func (kc *KafkaConsumer) Start() {
	kc.isRunning = true
	kc.readMessageLoop()
}

func (kc *KafkaConsumer) readMessageLoop() {

	for kc.isRunning {
		msg, err := kc.consumer.ReadMessage(-1)
		if err != nil {
			logrus.Errorf("Error reading message: %s", err)
			continue
		}

		var data types.ObuData
		if err := json.Unmarshal(msg.Value, &data); err != nil {
			logrus.Errorf("JSON serilization error: %s", err)
			continue
		}

		distance, err := kc.calcService.CalculateDistance(data)
		if err != nil {
			logrus.Errorf("distance calculation error: %s", err)
			continue
		}

		fmt.Printf("distance: %.2f \n", distance)
	}

}
