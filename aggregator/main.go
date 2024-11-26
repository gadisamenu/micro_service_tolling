package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"

	"github.com/gadisamenu/tolling/types"
)

func main() {
	listenAddr := flag.String("listenAddr", ":3000", "listen address of http server")

	store := NewMemoryStore()

	srvc := NewInvoiceAggregator(store)
	srvc = NewLogMiddleware(srvc)

	makeHTTPTransport(*listenAddr, srvc)

}
func makeHTTPTransport(listenAddr string, srvc Aggregator) {
	fmt.Println("Http listening on port: ", listenAddr)
	http.HandleFunc("/aggregate", handleAggregate(srvc))
	http.ListenAndServe(listenAddr, nil)

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
