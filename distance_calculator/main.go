package main

import (
	"log"
)

const topic = "obuTopic"

func main() {
	calcService := NewCalculatorService()
	calcService = NewLogMiddleware(calcService)
	kafkaConsumer, err := NewKafkaConsumer(topic, calcService)

	if err != nil {
		log.Fatal(err)
	}

	kafkaConsumer.Start()

}
