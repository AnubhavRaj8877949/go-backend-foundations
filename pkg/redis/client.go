package redis

import (
	"context"
	"fmt"

	"example.com/go-backend-foundation/internal/config"
	goredis "github.com/redis/go-redis/v9"
)

func NewClient(ctx context.Context, cfg config.RedisConfig) (*goredis.Client, error) {
	client := goredis.NewClient(&goredis.Options{
		Addr:     cfg.Address,
		Password: cfg.Password,
		DB:       cfg.Database,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		client.Close()
		return nil, fmt.Errorf("redis ping failed: %w", err)
	}

	return client, nil
}
