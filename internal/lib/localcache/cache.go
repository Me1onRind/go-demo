package localcache

import (
	"sync"

	"github.com/Me1onRind/go-demo/internal/dao/periodic_task_dao"
	"github.com/Me1onRind/go-demo/internal/lib/ctm_context"
)

func LoadCache(ctx *ctm_context.Context) {
	loaders := []Loader{
		periodic_task_dao.NewPeriodicTaskDao(),
	}
	wg := &sync.WaitGroup{}
	for _, loader := range loaders {
		wg.Add(1)
		go func(loader Loader) {
			defer wg.Done()
			if err := loadCache(ctx, loader); err != nil {
				panic(err)
			}
		}(loader)
	}

	listenCacheChange(ctx, loaders)
}
