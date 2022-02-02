package local_cache

import (
	"context"
	"fmt"
	"time"

	"github.com/Me1onRind/go-demo/global/client_singleton"
	"github.com/Me1onRind/go-demo/infrastructure/logger"
	"github.com/bluele/gcache"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

type Loader interface {
	DumpData() ([]*CacheItem, error)
}

type CacheItem struct {
	Namespace string
	Data      interface{}
}

type LocalCache struct {
	cache           gcache.Cache
	loaders         []Loader
	loaderInfoMap   map[string]*LoaderInfo
	refreshInterval time.Duration
}

type LoaderInfo struct {
	Loader  Loader
	Version int64
}

func NewLocalCache(loaders []Loader) *LocalCache {
	l := &LocalCache{
		cache:           gcache.New(len(loaders) * 10).Build(),
		loaders:         loaders,
		loaderInfoMap:   map[string]*LoaderInfo{},
		refreshInterval: time.Second * 3,
	}
	return l
}

func (l *LocalCache) InitLoad(ctx context.Context) {
	for _, loader := range l.loaders {
		err := l.load(ctx, loader)
		if err != nil {
			continue
		}
	}
}

func (l *LocalCache) Listen(ctx context.Context) {
	for {
		for namespace, info := range l.loaderInfoMap {
			version, err := getNamespaceVersion(ctx, namespace)
			if err != nil {
				continue
			}
			if info.Version < version {
				loaderInfo := l.loaderInfoMap[namespace]
				if loaderInfo == nil {
					logger.CtxError(ctx, "Empty loader info", zap.String("namespace", namespace))
					continue
				}

				if err := l.load(ctx, loaderInfo.Loader); err != nil {
					continue
				}
			}
			info.Version = version
		}
		time.Sleep(l.refreshInterval)
	}
}

func (l *LocalCache) Refresh(ctx context.Context, namespace string) error {
	version, err := client_singleton.RedisClient.Incr(ctx, getVersionKey(namespace)).Result()
	if err != nil {
		logger.CtxError(ctx, "LocalCache.Refresh fail", zap.String(namespace, namespace), zap.Error(err))
		return err
	}
	logger.CtxInfo(ctx, "LocalCache.Refresh", zap.String("namespace", namespace), zap.Int64("version", version))
	return nil
}

func (l *LocalCache) Query(ctx context.Context, namespace string) (interface{}, error) {
	m, err := l.cache.Get(namespace)
	if err != nil {
		logger.CtxError(ctx, "LocalCache.Query fail", zap.Error(err))
		return nil, err
	}
	return m, nil
}

func (l *LocalCache) load(ctx context.Context, loader Loader) error {
	itmes, err := loader.DumpData()
	if err != nil {
		logger.CtxError(ctx, "LocalCache.LoadAll fail", zap.Error(err))
		return err
	}
	for _, v := range itmes {
		l.cache.Set(v.Namespace, v.Data)
		if _, ok := l.loaderInfoMap[v.Namespace]; !ok {
			version, _ := getNamespaceVersion(ctx, v.Namespace)
			l.loaderInfoMap[v.Namespace] = &LoaderInfo{
				Loader:  loader,
				Version: version,
			}
		}
	}
	return nil
}

func getNamespaceVersion(ctx context.Context, namespace string) (int64, error) {
	version, err := client_singleton.RedisClient.Get(ctx, getVersionKey(namespace)).Int64()
	if err != nil {
		if err == redis.Nil {
			return 0, nil
		}
		logger.CtxError(ctx, "Get local cache version fail", zap.String("namespace", namespace), zap.Error(err))
	}
	return version, err
}

func getVersionKey(namespace string) string {
	return fmt.Sprintf("local:cache:%s", namespace)
}
