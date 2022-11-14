package logger

import (
	"context"
	"os"

	"github.com/sirupsen/logrus"
)

type loggerKey string

const (
	requestIdKey    loggerKey = "request-id"
	customPrefixKey loggerKey = "custom-prefix"
)

var (
	log = logrus.New()
)

func init() {
	log.Out = os.Stdout
}

func WithRequestId(ctx context.Context, requestId string) context.Context {
	return context.WithValue(ctx, requestIdKey, requestId)
}

func WithCustomPrefix(ctx context.Context, val any) context.Context {
	return context.WithValue(ctx, customPrefixKey, val)
}

func CtxDebugf(ctx context.Context, format string, a ...interface{}) {
	ctxLogf(ctx, logrus.DebugLevel, format, a...)
}

func Debugf(format string, a ...interface{}) {
	log.Debugf(format, a...)
}

func CtxInfof(ctx context.Context, format string, a ...interface{}) {
	ctxLogf(ctx, logrus.InfoLevel, format, a...)
}

func Infof(format string, a ...interface{}) {
	log.Infof(format, a...)
}

func CtxWarnf(ctx context.Context, format string, a ...interface{}) {
	ctxLogf(ctx, logrus.WarnLevel, format, a...)
}

func Warnf(format string, a ...interface{}) {
	log.Warnf(format, a...)
}

func CtxErrorf(ctx context.Context, format string, a ...interface{}) {
	ctxLogf(ctx, logrus.ErrorLevel, format, a...)
}

func Errorf(format string, a ...interface{}) {
	log.Infof(format, a...)
}

func ctxLogf(ctx context.Context, level logrus.Level, format string, a ...interface{}) {
	fields := logrus.Fields{}

	if reqId := ctx.Value(requestIdKey); reqId != nil {
		fields["request_id"] = reqId
	}

	if customPrefix := ctx.Value(customPrefixKey); customPrefix != nil {
		fields["custom_prefix"] = customPrefix
	}
	if len(fields) > 0 {
		log.WithFields(fields).Logf(level, format, a...)
		return
	}
	log.Logf(level, format, a...)
}
