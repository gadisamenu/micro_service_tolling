package main

import (
	"context"
	"encoding/json"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/gadisamenu/tolling/aggregator/client"
	"github.com/gadisamenu/tolling/types"
	"github.com/sirupsen/logrus"
)

type KafkaConsumer struct {
	consumer    *kafka.Consumer
	isRunning   bool
	calcService CalculatorServicer
	aggClient   client.Client
}

func NewKafkaConsumer(topic string, service CalculatorServicer, aggClient client.Client) (*KafkaConsumer, error) {
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
		aggClient:   aggClient,
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

		req := &types.AggregateRequest{
			Value: distance,
			Unix:  time.Now().Unix(),
			ObuId: int64(data.ObuId),
		}

		if err := kc.aggClient.Aggregate(context.Background(), req); err != nil {
			logrus.Error("aggregate error: ", err)
			continue
		}

	}

}
