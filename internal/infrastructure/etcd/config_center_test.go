package etcd

import (
	"context"
	"testing"
	"time"

	"github.com/Me1onRind/go-demo/internal/global/config"
	"github.com/Me1onRind/go-demo/internal/model/configmd"
	"github.com/stretchr/testify/assert"
	clientv3 "go.etcd.io/etcd/client/v3"
)

var client *clientv3.Client

func TestMain(m *testing.M) {
}

func Test_PullAndWatchConfig(t *testing.T) {
	config.LocalFileCfg = &configmd.LocalFileConfig{
		Etcd: configmd.EtcdConfig{
			Endpoints:   []string{"127.0.0.1:2379"},
			DialTimeout: time.Second * 5,
		},
	}
	err := PullAndWatchConfig(context.Background(), "/test/t2", nil)
	assert.ErrorContains(t, err, "context deadline exceeded")
	t.Log("test", err)
}
