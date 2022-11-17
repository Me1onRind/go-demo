package etcd

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	clientv3 "go.etcd.io/etcd/client/v3"
)

var (
	client *etcdClient
)

func TestMain(m *testing.M) {
	var err error
	localEtcdAddress := "127.0.0.1:2379"

	client, err = newEtcdClient(&clientv3.Config{
		Endpoints:   []string{localEtcdAddress},
		DialTimeout: time.Second * 2,
	})
	if err != nil {
		panic(err)
	}
	os.Exit(m.Run())
}

func Test_Put_Get(t *testing.T) {
	ctx := context.Background()
	key := "/unit_test_key"
	setValue := time.Now().String()
	timeout := time.Second * 2

	err := client.Put(ctx, key, setValue, timeout)
	if !assert.Empty(t, err) {
		return
	}

	value, err := client.Get(ctx, key, timeout)
	if !assert.Empty(t, err) {
		return
	}
	t.Logf("get value:[%s]", value)

	assert.Equal(t, setValue, string(value))
}
