package main

import (
	"context"

	"github.com/gadisamenu/tolling/types"
)

type GRPCAggregatorServer struct {
	types.UnimplementedAggregatorServer
	srvc Aggregator
}

func NewGRPCAggregatorServer(srvc Aggregator) *GRPCAggregatorServer {
	return &GRPCAggregatorServer{
		srvc: srvc,
	}
}

func (gr *GRPCAggregatorServer) Aggregate(ctx context.Context, aggReq *types.AggregateRequest) (*types.None, error) {

	distance := types.Distance{
		ObuId: int(aggReq.ObuId),
		Value: aggReq.Value,
		Unix:  aggReq.Unix,
	}

	return &types.None{}, gr.srvc.AggregateDistance(distance)

}
