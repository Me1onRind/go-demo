package localcache

import (
	"sync"
	"time"

	"github.com/Me1onRind/go-demo/internal/core/config"
	"github.com/Me1onRind/go-demo/internal/core/logger"
	"github.com/Me1onRind/go-demo/internal/lib/ctm_context"
	"github.com/Me1onRind/go-demo/internal/utils/goroutine"
	"github.com/bluele/gcache"
	"go.uber.org/zap"
)

var (
	versionMutex       = &sync.Mutex{}
	localcacheVersions = map[string]uint64{}
)

func listenCacheChange(ctx *ctm_context.Context, cache gcache.Cache, loaders []Loader) {
	goroutine.Go(func() {
		for {
			for _, loader := range loaders {
				localCacheKey := loader.GetLocalcacheKey()
				version, err := loader.LocalcacheVersion(ctx)
				if err != nil {
					continue
				}

				versionMutex.Lock()
				oldVersion, ok := localcacheVersions[localCacheKey]
				if !ok {
					localcacheVersions[localCacheKey] = version
					continue
				}
				versionMutex.Unlock()

				if version != oldVersion {
					logger.StdoutLogger.Info("Localcache version changed", zap.String("localCacheKey", localCacheKey),
						zap.Uint64("oldVersion", oldVersion), zap.Uint64("newVersion", version))
					goroutine.Go(func() {
						if err := loadCache(ctx, cache, loader); err == nil {
							versionMutex.Lock()
							localcacheVersions[localCacheKey] = version
							versionMutex.Unlock()
						}
					})
				}
			}

			time.Sleep(config.RemoteConfig.Localcache.GetCheckVersionInterval())
		}
	})
}

func loadCache(ctx *ctm_context.Context, cache gcache.Cache, loader Loader) error {
	beginTime := time.Now()
	defer func() {
		logger.StdoutLogger.Info("Load localcache", zap.String("localcacheKey", loader.GetLocalcacheKey()), zap.Duration("cost", time.Since(beginTime)))
	}()
	data, err := loader.LoadLocalcacheData(ctx)
	if err != nil {
		logger.StdoutLogger.Error("Load localcache failed", zap.String("error", err.String()))
		return err.GenError()
	}

	if err := cache.Set(loader.GetLocalcacheKey(), data); err != nil {
		logger.StdoutLogger.Error("Set localcache failed", zap.String("error", err.Error()))
		return err
	}
	return nil
}
