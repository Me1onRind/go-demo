package local_cache

import (
	"context"
	"testing"
	"time"

	"github.com/Me1onRind/go-demo/infrastructure/util/unit_test"
	"github.com/stretchr/testify/assert"
)

type testLoader struct {
	i int
}

func (t *testLoader) DumpData() ([]*CacheItem, error) {
	value := "value"
	if t.i > 0 {
		value = "value2"
	}
	t.i++
	return []*CacheItem{
		{
			Namespace: "ns",
			Data: map[string]string{
				"key": value,
			},
		},
	}, nil
}

func TestMain(m *testing.M) {
	unit_test.InitGlobalVar()
	m.Run()
}

func TestInitLoad(t *testing.T) {
	ctx := context.Background()
	l := NewLocalCache([]Loader{
		&testLoader{},
	})
	l.refreshInterval = time.Millisecond * 10

	t.Run("InitLoad", func(t *testing.T) {
		l.InitLoad(ctx)
		v, err := l.Query(ctx, "ns")
		if assert.Empty(t, err) {
			t.Log(v)
			data := v.(map[string]string)
			assert.Equal(t, "value", data["key"])
		}
	})

	go l.Listen(ctx)
	time.Sleep(time.Millisecond * 20)

	t.Run("Refresh", func(t *testing.T) {
		l.Refresh(ctx, "ns")
		time.Sleep(time.Millisecond * 100)
		v, err := l.Query(ctx, "ns")
		if assert.Empty(t, err) {
			t.Log(v)
			data := v.(map[string]string)
			assert.Equal(t, "value2", data["key"])
		}
	})
}
