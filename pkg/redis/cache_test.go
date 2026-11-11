package redis

import (
	"context"
	"testing"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

func TestSetAndGet(t *testing.T) {
	client := goredis.NewClient(&goredis.Options{
		Addr: "localhost:6379",
	})
	cache := NewCache(client)

	ctx := context.Background()
	key := "test:message"
	value := "hello redis"

	defer client.Close()
	defer cache.Delete(ctx, key)

	if err := client.Ping(ctx).Err(); err != nil {
		t.Fatalf("redis is not available: %v", err)
	}

	if err := cache.Set(ctx, key, value, time.Minute); err != nil {
		t.Fatalf("failed to set value: %v", err)
	}

	got, err := cache.Get(ctx, key)
	if err != nil {
		t.Fatalf("failed to get value: %v", err)
	}

	if got != value {
		t.Fatalf("expected %q, got %q", value, got)
	}

	if err := cache.Delete(ctx, key); err != nil {
		t.Fatalf("failed to delete value: %v", err)
	}

	_, err = cache.Get(ctx, key)
	if err != goredis.Nil {
		t.Fatalf("expected redis.Nil after delete, got %v", err)
	}
}

type testUser struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func TestSetAndGetJSON(t *testing.T) {
	client := goredis.NewClient(&goredis.Options{
		Addr: "localhost:6379",
	})
	cache := NewCache(client)

	ctx := context.Background()
	key := "test:user"

	defer client.Close()
	defer cache.Delete(ctx, key)

	if err := client.Ping(ctx).Err(); err != nil {
		t.Fatalf("redis is not available: %v", err)
	}

	expected := testUser{
		ID:    "1",
		Name:  "John",
		Email: "john@example.com",
	}

	if err := cache.SetJSON(ctx, key, expected, time.Minute); err != nil {
		t.Fatalf("failed to set JSON value: %v", err)
	}

	var actual testUser

	if err := cache.GetJSON(ctx, key, &actual); err != nil {
		t.Fatalf("failed to get JSON value: %v", err)
	}

	if actual != expected {
		t.Fatalf("expected %+v, got %+v", expected, actual)
	}
}
