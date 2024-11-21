package main

import (
	"math"

	"github.com/gadisamenu/tolling/types"
)

type CalculatorServicer interface {
	CalculateDistance(types.ObuData) (float64, error)
}

type CalculatorService struct {
	prevPoint []float64
}

func NewCalculatorService() CalculatorServicer {
	return &CalculatorService{}
}

func (cs *CalculatorService) CalculateDistance(data types.ObuData) (float64, error) {
	distance := 0.0
	if len(cs.prevPoint) > 0 {
		distance = calculateDistance(cs.prevPoint[0], cs.prevPoint[1], data.Lat, data.Long)
	}
	cs.prevPoint = []float64{data.Lat, data.Long}
	return distance, nil
}

func calculateDistance(x1, y1, x2, y2 float64) float64 {
	return math.Sqrt(math.Pow(x2-x1, 2) + math.Pow(y2-y1, 2))
}
