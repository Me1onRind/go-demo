package redis_client

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/alicebob/miniredis/v2/server"

	//"github.com/alicebob/miniredis/v2/server"
	"github.com/stretchr/testify/assert"
)

func Test_Set_Get(t *testing.T) {
	s, err := miniredis.Run()
	s.Server().SetPreHook(func(*server.Peer, string, ...string) bool {
		time.Sleep(time.Millisecond * 5)
		return false
	})
	assert.Empty(t, err)
	defer s.Close()
	r := NewRedisClientFromAddr(s.Addr())
	ctx := context.Background()
	err = r.Set(ctx, "test:key", "value", 0).Err()
	assert.Empty(t, err)
	val, err := r.Get(ctx, "test:key").Result()
	if assert.Empty(t, err) {
		t.Log(val)
		assert.Equal(t, "value", val)
	}
}
