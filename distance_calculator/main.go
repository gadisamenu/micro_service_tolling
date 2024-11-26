package main

import (
	"log"

	"github.com/gadisamenu/tolling/aggregator/client"
)

const (
	topic    = "obuTopic"
	endpoint = "http://127.0.0.1:3000/aggregate"
)

func main() {
	calcService := NewCalculatorService()
	calcService = NewLogMiddleware(calcService)
	aggClient := client.NewClient(endpoint)
	kafkaConsumer, err := NewKafkaConsumer(topic, calcService, aggClient)

	if err != nil {
		log.Fatal(err)
	}

	kafkaConsumer.Start()

}
