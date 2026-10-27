package service

import "context"

type (
	Currency interface {
		GetRate(ctx context.Context, from, to string) (float64, error)
		Convert(ctx context.Context, amount float64, from, to string) (float64, error)
	}

	Delivery interface {
		Calculate(ctx context.Context, currency string) (int64, error)
	}
)
