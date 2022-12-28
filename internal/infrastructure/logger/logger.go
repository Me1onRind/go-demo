package logger

import (
	"context"
	"os"

	"github.com/Me1onRind/go-demo/internal/infrastructure/tool/timehelper"
	"github.com/sirupsen/logrus"
)

type loggerKey struct{}

var (
	log = logrus.New()
)

func init() {
	log.Out = os.Stdout
	log.SetFormatter(&logrus.TextFormatter{
		DisableQuote:    true,
		TimestampFormat: timehelper.NormalFormat,
	})
	/*log.SetReportCaller(true)*/
}

func WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, loggerKey{}, log.WithContext(ctx))
}

func WithFields(ctx context.Context, fields logrus.Fields) context.Context {
	l := getLoggerFromCtx(ctx).WithFields(fields)
	return context.WithValue(ctx, loggerKey{}, l)
}

func getLoggerFromCtx(ctx context.Context) *logrus.Entry {
	if l, ok := ctx.Value(loggerKey{}).(*logrus.Entry); ok {
		return l
	}
	return log.WithContext(ctx)
}

func CtxDebugf(ctx context.Context, format string, a ...interface{}) {
	getLoggerFromCtx(ctx).Debugf(format, a...)
}

func Debugf(format string, a ...interface{}) {
	log.Debugf(format, a...)
}

func CtxInfof(ctx context.Context, format string, a ...interface{}) {
	getLoggerFromCtx(ctx).Infof(format, a...)
}

func Infof(format string, a ...interface{}) {
	log.Infof(format, a...)
}

func CtxWarnf(ctx context.Context, format string, a ...interface{}) {
	getLoggerFromCtx(ctx).Warnf(format, a...)
}

func Warnf(format string, a ...interface{}) {
	log.Warnf(format, a...)
}

func CtxErrorf(ctx context.Context, format string, a ...interface{}) {
	getLoggerFromCtx(ctx).Errorf(format, a...)
}

func Errorf(format string, a ...interface{}) {
	log.Errorf(format, a...)
}

func CtxFatalf(ctx context.Context, format string, a ...interface{}) {
	getLoggerFromCtx(ctx).Fatalf(format, a...)
}

func Fatalf(ctx context.Context, format string, a ...interface{}) {
	log.Fatalf(format, a...)
}
