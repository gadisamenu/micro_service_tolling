package main

import (
	"time"

	"github.com/gadisamenu/tolling/types"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/sirupsen/logrus"
)

type MetricsMiddleware struct {
	reqCounterAgg     prometheus.Counter
	reqCounterCalc    prometheus.Counter
	errReqCounterAgg  prometheus.Counter
	errReqCounterCalc prometheus.Counter
	reqLatencyAgg     prometheus.Histogram
	reqLatencyCalc    prometheus.Histogram
	next              Aggregator
}

func NewMetricsMiddleware(aggregator Aggregator) *MetricsMiddleware {
	reqCounterAgg := promauto.NewCounter(
		prometheus.CounterOpts{
			Namespace: "aggregator_request_counter",
			Name:      "aggregate",
		},
	)
	reqCounterCalc := promauto.NewCounter(
		prometheus.CounterOpts{
			Namespace: "aggregator_request_counter",
			Name:      "calculate",
		},
	)
	errReqCounterAgg := promauto.NewCounter(
		prometheus.CounterOpts{
			Namespace: "aggregator_request_counter_err",
			Name:      "aggregate",
		},
	)
	errReqCounterCalc := promauto.NewCounter(
		prometheus.CounterOpts{
			Namespace: "aggregator_request_counter_err",
			Name:      "calculate",
		},
	)

	reqLatencyAgg := promauto.NewHistogram(prometheus.HistogramOpts{
		Namespace: "aggregate_request_latecy",
		Name:      "aggregate",
		Buckets:   []float64{0.1, 0.5, 1.0},
	})

	reqLatencyCalc := promauto.NewHistogram(prometheus.HistogramOpts{
		Namespace: "aggregate_request_latecy",
		Name:      "calculate",
		Buckets:   []float64{0.1, 0.5, 1.0},
	})
	return &MetricsMiddleware{
		reqCounterAgg:     reqCounterAgg,
		reqCounterCalc:    reqCounterCalc,
		errReqCounterAgg:  errReqCounterAgg,
		errReqCounterCalc: errReqCounterCalc,
		reqLatencyAgg:     reqLatencyAgg,
		reqLatencyCalc:    reqLatencyCalc,
		next:              aggregator,
	}
}

func (mt *MetricsMiddleware) AggregateDistance(distance types.Distance) (err error) {
	defer func(start time.Time) {
		mt.reqLatencyAgg.Observe(float64(time.Since(start)))
		mt.reqCounterAgg.Inc()
		if err != nil {
			mt.errReqCounterAgg.Inc()
		}
	}(time.Now())
	err = mt.next.AggregateDistance(distance)
	return err
}

func (mt *MetricsMiddleware) CalculateInvoice(obuId int) (inv *types.Invoice, err error) {
	defer func(start time.Time) {
		mt.reqLatencyCalc.Observe(float64(time.Since(start)))
		mt.reqCounterCalc.Inc()
		if err != nil {
			mt.errReqCounterCalc.Inc()
		}
	}(time.Now())
	inv, err = mt.next.CalculateInvoice(obuId)
	return inv, err
}

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
