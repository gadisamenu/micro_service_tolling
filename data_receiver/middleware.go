package main

import (
	"time"

	"github.com/gadisamenu/tolling/types"
	"github.com/sirupsen/logrus"
)

type LogMiddleware struct {
	next DataProducer
}

func NewLogMiddleware(data DataProducer) *LogMiddleware {
	return &LogMiddleware{
		next: data,
	}
}

func (lg *LogMiddleware) ProduceData(data types.ObuData) error {
	go func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"obuId": data.ObuId,
			"lat":   data.Lat,
			"long":  data.Long,
			"time":  time.Since(start),
		}).Info("producing to kafka")

	}(time.Now())
	return lg.next.ProduceData(data)
}
