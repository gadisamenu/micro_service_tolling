package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gadisamenu/tolling/types"
)

type HTTPClient struct {
	Endpoint string
}

func NewHTTPClient(endpoint string) Client {
	return &HTTPClient{
		Endpoint: endpoint,
	}
}

func (c *HTTPClient) GetInvoice(ctx context.Context, id int) (*types.Invoice, error) {

	b, err := json.Marshal(id)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", c.Endpoint+"/invoice?obu="+strconv.Itoa(id), bytes.NewReader(b))

	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("service responded with non %d status code %d", http.StatusOK, res.StatusCode)
	}

	var inv types.Invoice

	if err := json.NewDecoder(res.Body).Decode(&inv); err != nil {
		return nil, err
	}
	defer res.Body.Close()
	return &inv, nil
}

func (c *HTTPClient) Aggregate(ctx context.Context, distance *types.AggregateRequest) error {

	b, err := json.Marshal(distance)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", c.Endpoint+"/aggregate", bytes.NewReader(b))
	if err != nil {
		return err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("service responded with non %d status code %d", http.StatusOK, res.StatusCode)
	}
	res.Body.Close()
	return nil

}
