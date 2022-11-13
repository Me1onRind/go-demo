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

func CtxInfof(ctx context.Context, format string, a ...interface{}) {
	fields := logrus.Fields{}

	if reqId := ctx.Value(requestIdKey); reqId != nil {
		fields["request_id"] = reqId
	}

	if customPrefix := ctx.Value(customPrefixKey); customPrefix != nil {
		fields["custom_prefix"] = customPrefix
	}

	if len(fields) > 0 {
		log.WithFields(fields).Infof(format, a...)
		return
	}

	log.Infof(format, a...)
}

func Infof(format string, a ...interface{}) {
	log.Infof(format, a...)
}
