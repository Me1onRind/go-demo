package middleware

import (
	"strings"

	"github.com/Me1onRind/go-demo/internal/infrastructure/logger"
	"github.com/Me1onRind/go-demo/internal/infrastructure/tracer"
	"github.com/gin-gonic/gin"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"go.elastic.co/apm/module/apmhttp"
)

func Tracer() gin.HandlerFunc {
	return func(c *gin.Context) {
		//span := opentracing.GlobalTracer().Extract()
		span := opentracing.GlobalTracer().StartSpan(c.Request.URL.Path)
		defer span.Finish()

		traceId, _ := traceIdAndSpanIdFromSpan(span.Context())
		setGinExtractContext(c, logger.WithFields(
			mustGetGinExtractContext(c),
			logrus.Fields{
				logger.TraceIdKey: traceId,
				//logger.SpanIdKey:  spanId,
			}))
		setGinExtractContext(c, tracer.WithSpan(mustGetGinExtractContext(c), span))

		c.Next()
	}
}

func traceIdAndSpanIdFromSpan(sm opentracing.SpanContext) (string, string) {
	carrier := opentracing.TextMapCarrier{}
	_ = opentracing.GlobalTracer().Inject(sm, opentracing.TextMap, &carrier)
	if v, ok := carrier[apmhttp.W3CTraceparentHeader]; ok {
		return traceIdAndSpanIdFromW3CTraceparent(v)
	}
	return "", ""
}

func traceIdAndSpanIdFromW3CTraceparent(value string) (string, string) {
	arr := strings.Split(value, "-")
	if len(arr) >= 3 {
		return arr[1], arr[2]
	}
	return "", ""
}
