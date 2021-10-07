package middleware

import (
	"context"
	"strings"

	"github.com/Me1onRind/go-demo/internal/lib/ctm_context"
	"github.com/gin-gonic/gin"
	opentracing "github.com/opentracing/opentracing-go"
	"go.elastic.co/apm/module/apmhttp"
	"go.uber.org/zap"
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

		commonCtx := ctx.(*ctm_context.Context)
		commonCtx.Span = span
		commonCtx.Logger = commonCtx.Logger.With(zap.String("requestID", requestIDFromSpan(span.Context())))

		resp, err = handler(ctx, req)
		return resp, err
	}
}

func GinTracer() gin.HandlerFunc {
	return func(c *gin.Context) {
		commonCtx := ctm_context.GetCtmCtxFromGinCtx(c)
		span := opentracing.GlobalTracer().StartSpan(c.Request.URL.Path)
		defer span.Finish()

		commonCtx.Span = span
		commonCtx.Logger = commonCtx.Logger.With(zap.String("requestID", requestIDFromSpan(span.Context())))

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
