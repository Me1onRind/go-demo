package logger

import (
	"context"
	"os"
	"time"

	"github.com/Me1onRind/go-demo/internal/infrastructure/tool/timehelper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type loggerKey struct{}

var (
	globalLogger *zap.SugaredLogger
)

func init() {
	encoder := zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
		MessageKey:       "msg",
		LevelKey:         "level",
		TimeKey:          "ts",
		CallerKey:        "file",
		EncodeCaller:     zapcore.ShortCallerEncoder,
		ConsoleSeparator: "|",
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format(timehelper.NormalFormat))
		},
		EncodeDuration: func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendInt64(int64(d) / 1000000)
		},
	})
	core := zapcore.NewTee(zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zap.DebugLevel))
	globalLogger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1),
		zap.Hooks(func(e zapcore.Entry) error {
			return nil
		})).
		Sugar()
}

func WithFields(ctx context.Context, args ...any) context.Context {
	l := getLoggerFromCtx(ctx).With(args...)
	return context.WithValue(ctx, loggerKey{}, l)
}

func getLoggerFromCtx(ctx context.Context) *zap.SugaredLogger {
	if l, ok := ctx.Value(loggerKey{}).(*zap.SugaredLogger); ok {
		return l
	}
	return globalLogger
}

func CtxDebugf(ctx context.Context, format string, a ...any) {
	getLoggerFromCtx(ctx).Debugf(format, a...)
}

func Debugf(format string, a ...any) {
	globalLogger.Debugf(format, a...)
}

func CtxInfof(ctx context.Context, format string, a ...any) {
	getLoggerFromCtx(ctx).Infof(format, a...)
}

func Infof(format string, a ...any) {
	globalLogger.Infof(format, a...)
}

func CtxWarnf(ctx context.Context, format string, a ...any) {
	getLoggerFromCtx(ctx).Warnf(format, a...)
}

func Warnf(format string, a ...any) {
	globalLogger.Warnf(format, a...)
}

func CtxErrorf(ctx context.Context, format string, a ...any) {
	getLoggerFromCtx(ctx).Errorf(format, a...)
}

func Errorf(format string, a ...any) {
	globalLogger.Errorf(format, a...)
}

func CtxFatalf(ctx context.Context, format string, a ...any) {
	getLoggerFromCtx(ctx).Fatalf(format, a...)
}

func Fatalf(ctx context.Context, format string, a ...any) {
	globalLogger.Fatalf(format, a...)
}
