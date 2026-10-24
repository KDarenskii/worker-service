package rcurrency

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/KDarenskii/worker-service/internal/app/entity"
	"github.com/KDarenskii/worker-service/internal/app/repository"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

const rateKeyPrefix = "rates:"

type repoRedis struct {
	client   *redis.Client
	cacheTTL time.Duration
}

func NewRepoFromRedis(client *redis.Client, cacheTTL time.Duration) repository.CurrencyRate {
	return &repoRedis{client: client, cacheTTL: cacheTTL}
}

func (r *repoRedis) GetRate(ctx context.Context, from, to string) (float64, error) {
	res, err := r.client.Get(ctx, r.buildKey(from, to)).Result()

	if errors.Is(err, redis.Nil) {
		return 0, entity.ErrRateNotFound
	}

	if err != nil {
		return 0, err
	}

	rate, err := strconv.ParseFloat(res, 64)
	if err != nil {
		return 0, fmt.Errorf("parse cached rate %q: %w", res, err)
	}

	return rate, nil
}

func (r *repoRedis) SetRate(ctx context.Context, from, to string, rate float64) error {
	strRate := strconv.FormatFloat(rate, 'f', -1, 64)

	if err := r.client.Set(ctx, r.buildKey(from, to), strRate, r.cacheTTL).Err(); err != nil {
		return fmt.Errorf("redis set rate %s/%s: %w", from, to, err)
	}

	return nil
}

func (r *repoRedis) SetRates(ctx context.Context, from string, rates map[string]float64) error {
	pipeliner := r.client.Pipeline()

	for to, value := range rates {
		pipeliner.Set(ctx, r.buildKey(from, to), value, r.cacheTTL)
	}

	_, err := pipeliner.Exec(ctx)
	if err != nil {
		return fmt.Errorf("redis set rates %s: %w", from, err)
	}

	log.Info().Ctx(ctx).Str("from", from).Int("count", len(rates)).Msg("Saved rates to redis cache")

	return nil
}

func (r *repoRedis) buildKey(from, to string) string {
	return fmt.Sprintf("%s%s:%s", rateKeyPrefix, from, to)
}
