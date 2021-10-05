package localcache

import (
	"context"
	"sync"

	"github.com/Me1onRind/go-demo/internal/dao/periodic_task"
)

func LoadConfigCache(ctx context.Context) {
	loaders := []Loader{
		periodic_task.NewPeriodicTaskDao(),
	}
	wg := &sync.WaitGroup{}
	for _, loader := range loaders {
		wg.Add(1)
		go func(loader Loader) {
			defer wg.Done()
			if err := loadCache(loader); err != nil {
				panic(err)
			}
		}(loader)
	}

	listenCacheChange(ctx, loaders)
}
