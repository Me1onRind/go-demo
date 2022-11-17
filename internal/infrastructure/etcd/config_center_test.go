package etcd

import (
	"context"
	"testing"
	"time"

	"github.com/Me1onRind/go-demo/internal/model/configmd"
)

func Test_PullAndWatchConfig(t *testing.T) {
	_ = PullAndWatchConfig(context.Background(), &configmd.EtcdConfig{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: time.Second * 5,
	}, "test", nil)
	t.Log("test")
}
