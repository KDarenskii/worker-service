package sdelivery

import (
	"context"
	"fmt"
	"math"

	"github.com/KDarenskii/worker-service/internal/app/service"
)

const (
	baseCurrency      = "EUR"
	baseDeliveryPrice = 1000
)

type srv struct {
	currency service.Currency
}

func NewService(currency service.Currency) service.Delivery {
	return &srv{currency: currency}
}

func (s *srv) Calculate(ctx context.Context, currency string) (int64, error) {
	price, err := s.currency.Convert(ctx, baseDeliveryPrice, baseCurrency, currency)
	if err != nil {
		return 0, fmt.Errorf("convert delivery price to %s: %w", currency, err)
	}

	return int64(math.Round(price)), nil
}
