package scurrency

import (
	"context"
	"errors"
	"fmt"

	"github.com/KDarenskii/worker-service/internal/app/entity"
	"github.com/KDarenskii/worker-service/internal/app/repository"
	"github.com/KDarenskii/worker-service/internal/app/service"
	"github.com/rs/zerolog/log"
)

type RatesProvider interface {
	GetRates(ctx context.Context, base string) (map[string]float64, error)
}

type srv struct {
	rates RatesProvider
	repo  repository.CurrencyRate
}

func NewService(rates RatesProvider, repo repository.CurrencyRate) service.Currency {
	return &srv{rates: rates, repo: repo}
}

func (s *srv) GetRate(ctx context.Context, from, to string) (float64, error) {
	if from == to {
		return 1, nil
	}

	if cacheRate, err := s.repo.GetRate(ctx, from, to); err == nil {
		return cacheRate, nil
	} else if !errors.Is(err, entity.ErrRateNotFound) {
		log.Warn().Ctx(ctx).Str("from", from).Str("to", to).Msg("Failed to get rate from cache")
	}

	rates, err := s.rates.GetRates(ctx, from)
	if err != nil {
		return 0, fmt.Errorf("get rates: %w", err)
	}

	err = s.repo.SetRates(ctx, from, rates)
	if err != nil {
		log.Error().Err(err).Str("from", from).Str("to", to).Msg("Failed to save rates in cache")
	}

	toRate, ok := rates[to]

	if !ok {
		return 0, entity.ErrFixerCurrencyNotFound
	}

	return toRate, nil
}

func (s *srv) Convert(ctx context.Context, amount float64, from, to string) (float64, error) {
	rate, err := s.GetRate(ctx, from, to)
	if err != nil {
		return 0, err
	}
	return amount * rate, nil
}
