package etcd

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Me1onRind/go-demo/internal/infrastructure/logger"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type Client interface {
	Get(ctx context.Context, key string, timeout time.Duration) ([]byte, error)
	Watch(ctx context.Context, key string, f WatchHandler)
	Put(ctx context.Context, key, value string, timeout time.Duration) error
	Close() error
}

var (
	ErrKeyNotFound = errors.New("Key not found")
)

type WatchHandler func(context.Context, *clientv3.Event)

type etcdClient struct {
	invoker *clientv3.Client
}

func NewEtcdClient(cfg *clientv3.Config) (Client, error) {
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

	logger.CtxInfof(ctx, "Put key:[%s], value:\n%s", key, value)
	_, err := e.invoker.Put(ctx, key, value)
	if err != nil {
		return fmt.Errorf("Put key:[%s] error, cause:[%w]", key, err)
	}
	logger.CtxInfof(ctx, "Put key:[%s] success", key)

	return nil
}

func (e *etcdClient) Watch(ctx context.Context, key string, f WatchHandler) {
	ch := e.invoker.Watch(ctx, key)
	logger.CtxInfof(ctx, "Start watch key:[%s]", key)
	select {
	case resp := <-ch:
		logger.CtxInfof(ctx, "Get watch respone key:[%s]", key)
		err := resp.Err()
		if err != nil {
			logger.CtxErrorf(ctx, "Watch key:[%s] fail", key)
			break
		}
		for _, event := range resp.Events {
			f(ctx, event)
		}
	case <-ctx.Done():
		return
	}
}

func (e *etcdClient) Close() error {
	return e.invoker.Close()
}
