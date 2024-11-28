package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"strconv"

	"github.com/gadisamenu/tolling/aggregator/client"
	"github.com/sirupsen/logrus"
)

type apiFunc func(http.ResponseWriter, *http.Request) error

func main() {

	listenAddr := flag.String("listenAddr", ":6000", "gate way http listening address")
	flag.Parse()

	client := client.NewHTTPClient("http://127.0.0.1:3000")
	invoiceHandler := NewInvoiceHandler(client)

	http.HandleFunc("/invoice/{obuId}", makeApiFunc(invoiceHandler.handleInvoice))

	logrus.Infof("getway http server running on port 6000")
	err := http.ListenAndServe(*listenAddr, nil)
	if err != nil {
		log.Fatal(err)
	}
}

type InvoiceHandler struct {
	client client.Client
}

func NewInvoiceHandler(client client.Client) *InvoiceHandler {
	return &InvoiceHandler{
		client: client,
	}
}

func (c *InvoiceHandler) handleInvoice(w http.ResponseWriter, r *http.Request) error {

	id, err := strconv.Atoi(r.PathValue("obuId"))
	if err != nil {
		return nil
	}

	inv, err := c.client.GetInvoice(r.Context(), id)
	if err != nil {
		return err
	}
	return writeJSON(w, http.StatusOK, inv)
}

func makeApiFunc(fn apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
	}
}

func writeJSON(w http.ResponseWriter, code int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(v)
}
