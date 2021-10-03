package localcache

import (
	"github.com/Me1onRind/go-demo/internal/core/logger"
	"github.com/Me1onRind/go-demo/internal/dao/periodic_task"
	"go.uber.org/zap"
)

func LoadConfigCache() {
	loaders := []Loader{
		periodic_task.NewPeriodicTaskDao(),
	}

	for _, loader := range loaders {
		if err := loader.LoadLocalCacheData(); err != nil {
			logger.Logger.Error("load localcache failed", zap.String("detail", err.String()))
		}
	}
}
