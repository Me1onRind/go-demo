package logger

import (
	"context"

	"github.com/Me1onRind/go-demo/constant/sys_constant"
	"github.com/Me1onRind/go-demo/global/logger_singleton"
	"go.uber.org/zap"
)

func CtxInfo(ctx context.Context, msg string, fields ...zap.Field) {
	value := ctx.Value(sys_constant.LoggerKey)
	if loggerInstantce, ok := value.(*zap.Logger); ok {
		loggerInstantce.Info(msg, fields...)
		return
	}
	logger_singleton.Logger.Info(msg, fields...)
}

func CtxWarn(ctx context.Context, msg string, fields ...zap.Field) {
	value := ctx.Value(sys_constant.LoggerKey)
	if loggerInstantce, ok := value.(*zap.Logger); ok {
		loggerInstantce.Warn(msg, fields...)
		return
	}
	logger_singleton.Logger.Warn(msg, fields...)
}

func CtxError(ctx context.Context, msg string, fields ...zap.Field) {
	value := ctx.Value(sys_constant.LoggerKey)
	if loggerInstantce, ok := value.(*zap.Logger); ok {
		loggerInstantce.Error(msg, fields...)
		return
	}
	logger_singleton.Logger.Error(msg, fields...)
}
