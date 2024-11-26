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
