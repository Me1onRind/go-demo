package etcd

import (
	"context"
	"errors"
	"fmt"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type Client interface {
	Get(ctx context.Context, key string, timeout time.Duration) ([]byte, error)
	Put(ctx context.Context, key, value string, timeout time.Duration) error
}

var (
	ErrKeyNotFound = errors.New("Key not found")
)

type etcdClient struct {
	invoker *clientv3.Client
}

func newEtcdClient(cfg *clientv3.Config) (*etcdClient, error) {
	invoker, err := clientv3.New(*cfg)
	if err != nil {
		return nil, err
	}
	e := &etcdClient{
		invoker: invoker,
	}
	return e, nil
}

func (e *etcdClient) Get(ctx context.Context, key string, timeout time.Duration) ([]byte, error) {
	if timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()
	}

	resp, err := e.invoker.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	if len(resp.Kvs) == 0 {
		return nil, fmt.Errorf("Get key:[%s] error, cause:[%w]", key, ErrKeyNotFound)
	}

	return resp.Kvs[0].Value, nil
}

func (e *etcdClient) Put(ctx context.Context, key, value string, timeout time.Duration) error {
	if timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()
	}

	_, err := e.invoker.Put(ctx, key, value)
	if err != nil {
		return fmt.Errorf("Put key:[%s] error, cause:[%w]", key, err)
	}

	return nil
}
