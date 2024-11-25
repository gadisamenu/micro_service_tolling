package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/gadisamenu/tolling/types"
	"github.com/gorilla/websocket"
)

var sendInterval = time.Second * 5

const wsEndpoint = "ws://127.0.0.1:30000/ws"

func genCoord() float64 {
	num := float64(rand.Intn(100) + 1)
	decimal := rand.Float64()

	return num + decimal
}

func genLatLong() (float64, float64) {
	return genCoord(), genCoord()
}

func genObuIds(n int) []int {
	obuIds := make([]int, n)
	for i := range n {
		obuIds[i] = rand.Intn(math.MaxInt)
	}
	return obuIds
}

func main() {
	obuIds := genObuIds(20)
	conn, _, err := websocket.DefaultDialer.Dial(wsEndpoint, nil)
	if err != nil {
		log.Fatal(err)
	}

	for {
		time.Sleep(sendInterval)
		for i := 0; i < len(obuIds); i++ {
			lat, long := genLatLong()
			data := types.ObuData{
				ObuId: obuIds[i],
				Lat:   lat,
				Long:  long,
			}
			if err := conn.WriteJSON(data); err != nil {
				fmt.Println(err)
			}

		}

	}
}
