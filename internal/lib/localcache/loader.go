package localcache

import (
	"context"
	"fmt"
	"time"

	"github.com/Me1onRind/go-demo/internal/core/logger"
	"github.com/Me1onRind/go-demo/internal/core/store"
	"github.com/Me1onRind/go-demo/internal/utils/goroutine"
	"go.uber.org/zap"
)

var (
	localcacheVersions = map[string]uint64{}
)

func listenCacheChange(ctx context.Context, loaders []Loader) {
	goroutine.Go(func() {
		for {
			for _, loader := range loaders {
				localCacheKey := loader.LocalCacheKey()
				vk := versionKey(localCacheKey)
				version, err := store.RedisPool.Get(ctx, vk).Uint64()
				if err != nil {
					logger.StdoutLogger.Error("Get localcache version failed", zap.String("localCacheKey", localCacheKey), zap.Error(err))
				}

				oldVersion, ok := localcacheVersions[localCacheKey]
				if !ok {
					localcacheVersions[localCacheKey] = version
					continue
				}

				if version != oldVersion {
					logger.StdoutLogger.Info("Localcache version changed", zap.String("localCacheKey", localCacheKey),
						zap.Uint64("oldVersion", oldVersion), zap.Uint64("newVersion", version))
					goroutine.Go(func() {
						_ = loadCache(loader)
					})
				}
			}
		}
	})
}

func loadCache(loader Loader) error {
	beginTime := time.Now()
	defer func() {
		logger.StdoutLogger.Info("Load localcache", zap.String("localCacheKey", loader.LocalCacheKey()), zap.Duration("cost", time.Since(beginTime)))
	}()
	if err := loader.LoadLocalCacheData(); err != nil {
		logger.StdoutLogger.Error("Load localcache failed", zap.String("error", err.String()))
		return err.GenError()
	}
	return nil
}

func versionKey(key string) string {
	return fmt.Sprintf("localcache:version:%s", key)
}
