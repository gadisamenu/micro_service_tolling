package main

import (
	"time"

	"github.com/gadisamenu/tolling/types"
	"github.com/sirupsen/logrus"
)

type LogMiddleware struct {
	next Aggregator
}

func NewLogMiddleware(aggregator Aggregator) Aggregator {
	return &LogMiddleware{
		next: aggregator,
	}
}

func (l *LogMiddleware) AggregateDistance(distance types.Distance) (err error) {

	defer func(start time.Time) {
		logrus.WithFields(
			logrus.Fields{
				"took": time.Since(start),
				"err":  err,
			},
		).Info("Aggregate distance")

	}(time.Now())

	err = l.next.AggregateDistance(distance)
	return err
}

func (l *LogMiddleware) CalculateInvoice(obuId int) (inv *types.Invoice, err error) {
	var (
		distance float64
		amount   float64
	)
	if inv != nil {
		distance = inv.TotalDistance
		amount = inv.TotalAmount
	}

	defer func(start time.Time) {
		logrus.WithFields(
			logrus.Fields{
				"took":          time.Since(start),
				"err":           err,
				"obuId":         obuId,
				"totalDistance": distance,
				"totalAmount":   amount,
			},
		).Info("Calculate invoice")

	}(time.Now())

	inv, err = l.next.CalculateInvoice(obuId)
	return inv, err
}
