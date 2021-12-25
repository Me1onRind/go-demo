package middleware

import (
	"context"
	"strings"

	"github.com/Me1onRind/go-demo/constant/sys_constant"
	"github.com/gin-gonic/gin"
	opentracing "github.com/opentracing/opentracing-go"
	"go.elastic.co/apm/module/apmhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func GrpcTracer() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		// tracer
		md, _ := metadata.FromIncomingContext(ctx)
		carrier := opentracing.HTTPHeadersCarrier{
			apmhttp.W3CTraceparentHeader: md.Get(strings.ToLower(apmhttp.W3CTraceparentHeader)),
			apmhttp.TracestateHeader:     md.Get(strings.ToLower(apmhttp.TracestateHeader)),
		}
		spanCtx, _ := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, carrier)
		var span opentracing.Span
		if spanCtx == nil {
			span = opentracing.GlobalTracer().StartSpan(info.FullMethod)
		} else {
			span = opentracing.GlobalTracer().StartSpan(info.FullMethod, opentracing.ChildOf(spanCtx))
		}
		defer func() {
			span.SetTag("req", req)
			span.SetTag("resp", resp)
			span.SetTag("error", err)
			span.Finish()
		}()

		ctx = context.WithValue(ctx, sys_constant.SpanKey, span)

		resp, err = handler(ctx, req)
		return resp, err
	}
}

func GinTracer() gin.HandlerFunc {
	return func(c *gin.Context) {
		span := opentracing.GlobalTracer().StartSpan(c.Request.URL.Path)
		defer span.Finish()

		traceId := requestIDFromSpan(span.Context())
		c.Set(sys_constant.TraceIdKey, traceId)
		c.Next()
	}
}

func requestIDFromSpan(sm opentracing.SpanContext) string {
	carrier := opentracing.TextMapCarrier{}
	_ = opentracing.GlobalTracer().Inject(sm, opentracing.TextMap, &carrier)
	if v, ok := carrier[apmhttp.W3CTraceparentHeader]; ok {
		return requestIDFromW3CTraceparent(v)
	}
	return ""
}

func requestIDFromW3CTraceparent(value string) string {
	arr := strings.Split(value, "-")
	if len(arr) >= 2 {
		return arr[1]
	}
	return ""
}
