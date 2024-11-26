package main

import (
	"github.com/gadisamenu/tolling/types"
)

const basePrice = 2.234

type Aggregator interface {
	AggregateDistance(types.Distance) error
	CalculateInvoice(int) (*types.Invoice, error)
}

type Storer interface {
	Insert(types.Distance) error
	Get(int) (float64, error)
}

type InvoiceAggregator struct {
	store Storer
}

func NewInvoiceAggregator(store Storer) Aggregator {
	return &InvoiceAggregator{
		store: store,
	}
}

func (i *InvoiceAggregator) AggregateDistance(distance types.Distance) error {
	return i.store.Insert(distance)
}

func (i *InvoiceAggregator) CalculateInvoice(obuId int) (*types.Invoice, error) {

	dist, err := i.store.Get(obuId)
	if err != nil {
		return nil, err
	}

	invoice := &types.Invoice{
		ObuId:         obuId,
		TotalDistance: dist,
		TotalAmount:   basePrice * dist,
	}
	return invoice, nil
}
