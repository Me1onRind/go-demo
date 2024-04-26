package redis

import (
	"context"
	"testing"
	"time"

	"github.com/Me1onRind/go-demo/internal/model/configmd"
	"github.com/alicebob/miniredis/v2"
	"github.com/alicebob/miniredis/v2/server"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

func Test_Base(t *testing.T) {
	redisServer := miniredis.RunT(t)
	redisServer.Server().SetPreHook(func(p *server.Peer, s string, s1 ...string) bool {
		time.Sleep(time.Millisecond * 2)
		if s == "GET" && s1[0] == "bad" {
			p.WriteError("mock error")
		}
		return false
	})
	defer redisServer.Close()
	cfg := configmd.RedisConfig{
		Addr: redisServer.Addr(),
	}
	client, err := NewRedisPool(&cfg)
	if !assert.Empty(t, err) {
		return
	}

	ctx := context.Background()
	client.Set(ctx, "k", "value", time.Second*10)
	value, err := client.Get(ctx, "k").Bytes()
	assert.Empty(t, err)
	assert.Equal(t, "value", string(value))
	client.Get(ctx, "bad")
}

func TestSetGetRedis(t *testing.T) {
	client := &redis.Client{}
	err := RegisterRedisClient("redis", client)
	assert.Empty(t, err)

	getClient := GetRedisClient("redis")
	assert.Equal(t, getClient, client)
}
