package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"net/http"
	"strconv"

	"github.com/gadisamenu/tolling/types"
	"google.golang.org/grpc"
)

func main() {
	httpListenAddr := flag.String("httpListenAddr", ":3000", "listen address of http server")
	grpcListenAddr := flag.String("grpcListenAddr", ":3001", "listen address of grpc server")

	store := NewMemoryStore()

	srvc := NewInvoiceAggregator(store)
	srvc = NewLogMiddleware(srvc)

	go makeGRPCTransport(*grpcListenAddr, srvc)
	makeHTTPTransport(*httpListenAddr, srvc)

}

func makeGRPCTransport(listenAddr string, srvc Aggregator) error {
	fmt.Println("GRPC listening on port: ", listenAddr)
	ln, err := net.Listen("TCP", listenAddr)
	if err != nil {
		return err
	}

	defer ln.Close()

	grpcServer := grpc.NewServer([]grpc.ServerOption{}...)

	types.RegisterAggregatorServer(grpcServer, NewGRPCAggregatorServer(srvc))

	return grpcServer.Serve(ln)
}
func makeHTTPTransport(listenAddr string, srvc Aggregator) {
	fmt.Println("Http listening on port: ", listenAddr)
	http.HandleFunc("/aggregate", handleAggregate(srvc))
	http.HandleFunc("/invoice", handleInvoice(srvc))
	http.ListenAndServe(listenAddr, nil)

}
func handleInvoice(srvc Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		values, ok := r.URL.Query()["obu"]
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
