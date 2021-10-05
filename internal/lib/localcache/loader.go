package localcache

import (
	"fmt"
	"time"

	"github.com/Me1onRind/go-demo/global/store"
	"github.com/Me1onRind/go-demo/internal/core/logger"
	"github.com/Me1onRind/go-demo/internal/lib/ctm_context"
	"github.com/Me1onRind/go-demo/internal/utils/goroutine"
	"go.uber.org/zap"
)

var (
	localcacheVersions = map[string]uint64{}
)

func listenCacheChange(ctx *ctm_context.Context, loaders []Loader) {
	goroutine.Go(func() {
		for {
			for _, loader := range loaders {
				localCacheKey := loader.LocalCacheKey()
				vk := versionKey(localCacheKey)
				version, err := store.RedisClient.Get(ctx, vk).Uint64()
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
						_ = loadCache(ctx, loader)
					})
				}
			}
		}
	})
}

func loadCache(ctx *ctm_context.Context, loader Loader) error {
	beginTime := time.Now()
	defer func() {
		logger.StdoutLogger.Info("Load localcache", zap.String("localCacheKey", loader.LocalCacheKey()), zap.Duration("cost", time.Since(beginTime)))
	}()
	data, err := loader.LoadLocalCacheData(ctx)
	if err != nil {
		logger.StdoutLogger.Error("Load localcache failed", zap.String("error", err.String()))
		return err.GenError()
	}

	if err := loader.LocalCacheInstance().Set(loader.LocalCacheKey(), data); err != nil {
		logger.StdoutLogger.Error("Set localcache failed", zap.String("error", err.Error()))
		return err
	}
	return nil
}

func versionKey(key string) string {
	return fmt.Sprintf("localcache:version:%s", key)
}
