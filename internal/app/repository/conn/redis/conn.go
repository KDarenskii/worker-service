package rcredis

import (
	"context"
	"fmt"
	"time"

	"github.com/KDarenskii/worker-service/internal/app/config/section"
	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

// Client встраивает *redis.Client — наружу доступны все команды go-redis напрямую.
type Client struct {
	*redis.Client
}

func NewClient(ctx context.Context, cfg section.RepositoryRedis) (*Client, error) {
	log.Info().Str("address", cfg.Address).Int("db_number", cfg.DB).Msg("redis connecting")

	redisClient := redis.NewClient(
		&redis.Options{
			Addr:     cfg.Address,
			Password: cfg.Password,
			DB:       cfg.DB,
		},
	)

	if err := redisotel.InstrumentTracing(redisClient); err != nil {
		closeClient(redisClient)
		return nil, fmt.Errorf("failed to initialize redis tracing: %w", err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	if err := redisClient.Ping(pingCtx).Err(); err != nil {
		closeClient(redisClient)
		return nil, fmt.Errorf("failed to ping redis: %w", err)
	}

	log.Info().Str("address", cfg.Address).Int("db_number", cfg.DB).Msg("redis connected")

	return &Client{redisClient}, nil
}

func closeClient(client *redis.Client) {
	if client == nil {
		return
	}

	if err := client.Close(); err != nil {
		log.Error().Err(err).Msg("redis: failed to close client")
	}
}
