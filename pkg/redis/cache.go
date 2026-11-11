package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

type Cache struct {
	client *goredis.Client
}

func NewCache(client *goredis.Client) *Cache {
	return &Cache{
		client: client,
	}
}
func (c *Cache) Set(
	ctx context.Context,
	key string,
	value string,
	expiration time.Duration,
) error {
	return c.client.Set(ctx, key, value, expiration).Err()
}
func (c *Cache) Get(
	ctx context.Context,
	key string,
) (string, error) {
	return c.client.Get(ctx, key).Result()
}

func (c *Cache) Delete(
	ctx context.Context,
	key string,
) error {
	return c.client.Del(ctx, key).Err()
}

func (c *Cache) SetJSON(
	ctx context.Context,
	key string,
	value any,
	expiration time.Duration,
) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("marshal cache value: %w", err)
	}

	return c.Set(ctx, key, string(data), expiration)
}

func (c *Cache) GetJSON(
	ctx context.Context,
	key string,
	destination any,
) error {
	data, err := c.Get(ctx, key)
	if err != nil {
		return fmt.Errorf("get cache value: %w", err)
	}

	if err := json.Unmarshal([]byte(data), destination); err != nil {
		return fmt.Errorf("unmarshal cache value: %w", err)
	}

	return nil
}
