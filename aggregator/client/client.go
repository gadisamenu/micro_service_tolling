package client

import (
	"context"

	"github.com/gadisamenu/tolling/types"
)

type Client interface {
	Aggregate(context.Context, *types.AggregateRequest) error
}
