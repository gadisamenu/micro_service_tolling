package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/gadisamenu/tolling/aggregator/client"
	"github.com/gadisamenu/tolling/types"
	"google.golang.org/grpc"
)

func main() {
	httpListenAddr := flag.String("httpListenAddr", ":3000", "listen address of http server")
	grpcListenAddr := flag.String("grpcListenAddr", ":3001", "listen address of grpc server")

	store := NewMemoryStore()

	srvc := NewInvoiceAggregator(store)
	srvc = NewLogMiddleware(srvc)

	go func() {
		log.Fatal(makeGRPCTransport(*grpcListenAddr, srvc))
	}()

	time.Sleep(time.Second * 2)

	client, err := client.NewGRPCClient(*grpcListenAddr)
	if err != nil {
		log.Fatal(err)
	}
	_ = client

	log.Fatal(makeHTTPTransport(*httpListenAddr, srvc))

}

func makeGRPCTransport(listenAddr string, srvc Aggregator) error {
	fmt.Println("GRPC listening on port: ", listenAddr)
	ln, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return err
	}

	defer ln.Close()

	grpcServer := grpc.NewServer([]grpc.ServerOption{}...)

	types.RegisterAggregatorServer(grpcServer, NewGRPCAggregatorServer(srvc))

	return grpcServer.Serve(ln)
}
func makeHTTPTransport(listenAddr string, srvc Aggregator) error {
	fmt.Println("Http listening on port: ", listenAddr)
	http.HandleFunc("/aggregate", handleAggregate(srvc))
	http.HandleFunc("/invoice", handleInvoice(srvc))
	return http.ListenAndServe(listenAddr, nil)

}
func handleInvoice(srvc Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.URL.Query())
		values, ok := r.URL.Query()["obu"]
		fmt.Println("obu", values, ok)
		if !ok {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "query obu id is required"})
			return
		}
		obuId, err := strconv.Atoi(values[0])
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid obu id"})
			return
		}

		inv, err := srvc.CalculateInvoice(obuId)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}

		writeJSON(w, http.StatusOK, inv)
	}
}

func handleAggregate(aggregator Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var distance types.Distance
		if err := json.NewDecoder(r.Body).Decode(&distance); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}

		if err := aggregator.AggregateDistance(distance); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}
