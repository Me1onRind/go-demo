package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/Me1onRind/go-demo/internal/infrastructure/logger"
	"github.com/Me1onRind/go-demo/internal/infrastructure/tool/random"
	"github.com/Me1onRind/go-demo/internal/infrastructure/tracer"
	"github.com/Me1onRind/go-demo/protocol/httpheader"
	"github.com/gin-gonic/gin"
	opentracing "github.com/opentracing/opentracing-go"
	"go.elastic.co/apm/module/apmhttp"
)

func Tracer(ctx context.Context, spanName string, header http.Header) (context.Context, opentracing.Span) {
	carrier := opentracing.HTTPHeadersCarrier(header)

	var traceId, spanId string
	var span opentracing.Span
	spanContext, err := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, carrier)

	if spanContext == nil {
		if err != nil && !errors.Is(err, opentracing.ErrSpanContextNotFound) {
			logger.Errorf("extract error:[%s], carrier:[%v]", err, carrier)
		}
		span = opentracing.GlobalTracer().StartSpan(spanName)
		carrier := opentracing.HTTPHeadersCarrier{}
		if err := opentracing.GlobalTracer().Inject(span.Context(), opentracing.HTTPHeaders, &carrier); err != nil {
			logger.Errorf("inject error:[%s], carrier:[%v]", err, carrier)
		} else {
			traceId, spanId = traceIdAndSpanIdFromSpan(carrier)
		}
	} else {
		span = opentracing.GlobalTracer().StartSpan(spanName, opentracing.ChildOf(spanContext))
		traceId, spanId = traceIdAndSpanIdFromSpan(carrier)
	}

	requestId := getRequestId(header)

	if len(traceId) == 0 || len(spanId) == 0 || len(requestId) == 0 {
		logger.Errorf("traceId:[%s],spanId:[%s],should be not happen,http_header:[%v]", traceId, spanId, carrier)
	}

	ctx = tracer.WithTracerId(ctx, traceId, spanId)
	ctx = tracer.WithSpan(ctx, span)
	ctx = tracer.WithRequestId(ctx, requestId)

	ctx = logger.WithFields(ctx, logger.RequestIdKey, requestId, logger.TraceIdKey, traceId, logger.SpanIdKey, spanId)
	return ctx, span

}

func GinTracer() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := mustGetGinExtractContext(c)
		ctx, span := Tracer(ctx, c.Request.URL.Path, c.Request.Header)
		defer span.Finish()
		setGinExtractContext(c, ctx)
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

func getRequestId(header http.Header) string {
	requestIdValaue := header[httpheader.RequestIdKey]
	if len(requestIdValaue) > 0 {
		return requestIdValaue[0]
	}
	return random.UUID()
}
