package store

import (
	"context"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
)

func Test_Set_Get(t *testing.T) {
	s, err := miniredis.Run()
	assert.Empty(t, err)
	defer s.Close()
	r := NewRedisPoolFromAddr(s.Addr())
	ctx := context.Background()
	err = r.Set(ctx, "test:key", "value", 0).Err()
	assert.Empty(t, err)
	val, err := r.Get(ctx, "test:key").Result()
	if assert.Empty(t, err) {
		t.Log(val)
		assert.Equal(t, "value", val)
	}
}
