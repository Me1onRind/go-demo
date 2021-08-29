package store

import (
	"context"
	"testing"

	"github.com/go-redis/redis/v8"
)

func Test_Set(t *testing.T) {
	r := NewRedisPool("localhost:6379")
	ctx := context.Background()
	err := r.Set(ctx, "test:key", "value", 0).Err()
	if err != nil {
		t.Fatal(err)
	}
}

func Test_Get(t *testing.T) {
	r := NewRedisPool("localhost:6379")
	ctx := context.Background()
	val, err := r.Get(ctx, "test:key").Result()
	if err != nil && err != redis.Nil {
		t.Fatal(err)
	}
	t.Log(val)
}
