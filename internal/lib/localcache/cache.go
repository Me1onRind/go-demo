package localcache

import (
	"sync"

	"github.com/Me1onRind/go-demo/internal/dao/kv_config_dao"
	"github.com/Me1onRind/go-demo/internal/lib/ctm_context"
	"github.com/bluele/gcache"
)

type cacheConfig struct {
	loaders []Loader
	cache   gcache.Cache
}

func newPermanentCacheConfig(loaders []Loader) *cacheConfig {
	return &cacheConfig{
		loaders: loaders,
		cache:   gcache.New(len(loaders)).Build(),
	}
}

func LoadCache(ctx *ctm_context.Context) {
	configs := []*cacheConfig{
		newPermanentCacheConfig([]Loader{
			kv_config_dao.NewKvConfigDao(),
		}),
	}

	wg := &sync.WaitGroup{}
	for _, config := range configs {
		for _, loader := range config.loaders {
			loader.SetLocalcache(config.cache)
			wg.Add(1)
			go func(loader Loader) {
				defer wg.Done()
				if err := loadCache(ctx, config.cache, loader); err != nil {
					panic(err)
				}
			}(loader)
		}
		listenCacheChange(ctx, config.cache, config.loaders)
	}
	wg.Wait()
}
