package goroutine

import (
	"context"
	"runtime/debug"

	"github.com/Me1onRind/go-demo/internal/infrastructure/logger"
)

func LogPanicStack(ctx context.Context, err any) {
	stack := debug.Stack()
	logger.CtxErrorf(ctx, "Panic:[%s]\n%s", err, stack)
}

func SafeGo(ctx context.Context, f func()) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				LogPanicStack(ctx, err)
			}
		}()
		f()
	}()
}
