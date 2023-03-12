package middleware

import (
	"errors"
	"strings"

	"github.com/Me1onRind/go-demo/internal/infrastructure/logger"
	"github.com/Me1onRind/go-demo/internal/infrastructure/tracer"
	"github.com/gin-gonic/gin"
	opentracing "github.com/opentracing/opentracing-go"
	"go.elastic.co/apm/module/apmhttp"
)

func Tracer() gin.HandlerFunc {
	return func(c *gin.Context) {
		var traceId, spanId string
		var span opentracing.Span

		spanContext, err := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(c.Request.Header))
		// root span
		if spanContext == nil {
			if err != nil && !errors.Is(err, opentracing.ErrSpanContextNotFound) {
				logger.Errorf("extract error:[%s], http_header:[%v]", err, c.Request.Header)
			}
			span = opentracing.GlobalTracer().StartSpan(c.Request.URL.Path)
			carrier := opentracing.HTTPHeadersCarrier{}
			if err := opentracing.GlobalTracer().Inject(span.Context(), opentracing.HTTPHeaders, &carrier); err != nil {
				logger.Errorf("inject error:[%s], http_header:[%v]", err, c.Request.Header)
			} else {
				traceId, spanId = traceIdAndSpanIdFromSpan(carrier)
			}
		} else {
			span = opentracing.GlobalTracer().StartSpan(c.Request.URL.Path, opentracing.ChildOf(spanContext))
			traceId, spanId = traceIdAndSpanIdFromSpan(opentracing.HTTPHeadersCarrier(c.Request.Header))
		}
		defer span.Finish()
		if len(traceId) == 0 || len(spanId) == 0 {
			logger.Errorf("traceId:[%s],spanId:[%s],should be not happen,http_header:[%v]", traceId, spanId, c.Request.Header)
		}

		setGinExtractContext(c, logger.WithFields(
			tracer.WithTracerId(mustGetGinExtractContext(c), traceId, spanId),
			logger.TraceIdKey, traceId,
		))
		setGinExtractContext(c, tracer.WithSpan(mustGetGinExtractContext(c), span))

		c.Next()
	}
}

func traceIdAndSpanIdFromSpan(carrier opentracing.HTTPHeadersCarrier) (string, string) {
	if v, ok := carrier[apmhttp.W3CTraceparentHeader]; ok {
		return traceIdAndSpanIdFromW3CTraceparent(v)
	}
	return "", ""
}

func traceIdAndSpanIdFromW3CTraceparent(values []string) (string, string) {
	if len(values) >= 1 {
		arr := strings.Split(values[0], "-")
		if len(arr) >= 3 {
			return arr[1], arr[2]
		}
	}
	return "", ""
}
