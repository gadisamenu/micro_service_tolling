package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gadisamenu/tolling/types"
	"github.com/gorilla/websocket"
)

var kafkaTopic = "obuTopic"

type DataReceiver struct {
	msgch    chan types.ObuData
	conn     *websocket.Conn
	producer DataProducer
}

func NewDataReceiver() (*DataReceiver, error) {

	kafka, err := NewKafkaProducer(kafkaTopic)

	logger := NewLogMiddleware(kafka)

	if err != nil {
		log.Fatal(err)
	}

	return &DataReceiver{
		msgch:    make(chan types.ObuData, 128),
		producer: logger,
	}, nil
}

func (dr *DataReceiver) wsHandler(w http.ResponseWriter, r *http.Request) {
	u := websocket.Upgrader{
		WriteBufferSize: 1028,
		ReadBufferSize:  1028,
	}

	conn, err := u.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	dr.conn = conn
	go dr.wsReceiveLoop()
}

func (dr *DataReceiver) wsReceiveLoop() {
	fmt.Println("client connected")
	for {
		var data types.ObuData
		if err := dr.conn.ReadJSON(&data); err != nil {
			fmt.Println("read error", err)
			continue
		}

		dr.producer.ProduceData(data)
	}

}

func main() {

	dr, err := NewDataReceiver()
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/ws", dr.wsHandler)
	http.ListenAndServe(":30000", nil)
}
