package goroutine

import (
	"runtime/debug"

	"github.com/Me1onRind/go-demo/internal/core/logger"
	"go.uber.org/zap"
)

func Go(f func()) {
	defer func() {
		if err := recover(); err != nil {
			logger.Logger.Error("Goroutine panic", zap.Any("panicErr", err))
			logger.Logger.Sugar().Errorf("%s", debug.Stack())
		}
	}()
	go f()
}
