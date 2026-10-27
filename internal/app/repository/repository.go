package repository

import "context"

type (
	CurrencyRate interface {
		GetRate(ctx context.Context, from, to string) (float64, error)
		SetRate(ctx context.Context, from, to string, rate float64) error
		SetRates(ctx context.Context, from string, rates map[string]float64) error
	}
)
