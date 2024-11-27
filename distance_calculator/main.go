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
	// httpClient := client.NewHTTPClient(endpoint)
	grpcClient, err := client.NewGRPCClient(":3001")
	if err != nil {
		log.Fatal(err)
	}
	kafkaConsumer, err := NewKafkaConsumer(topic, calcService, grpcClient)

	if err != nil {
		log.Fatal(err)
	}

	kafkaConsumer.Start()

}
