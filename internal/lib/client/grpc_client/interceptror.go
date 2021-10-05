package grpc_client

import (
	"context"
	"fmt"
	"time"

	"github.com/Me1onRind/go-demo/internal/lib/ctm_context"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type RetryFilter func(method string, req, reply interface{}, err error) bool

func withRetry(maxRetry uint32, needRetry RetryFilter) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		var err error
		for i := uint32(0); i < maxRetry+1; i++ {
			err = invoker(ctx, method, req, reply, cc, opts...)
			if err == nil {
				return nil
			}
			if needRetry != nil && !needRetry(method, req, reply, err) {
				return err
			}
		}
		return err
	}
}

func retryOnlyByCode(codes ...codes.Code) RetryFilter {
	return func(method string, req, reply interface{}, err error) bool {
		if s, ok := status.FromError(err); ok {
			for _, v := range codes {
				if s.Code() == v {
					return true
				}
			}
		}
		return false
	}
}

func withLogger() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		var err error
		begin := time.Now()
		defer func() {
			commonCtx := ctm_context.GetContext(ctx)
			commonCtx.Logger.Info("grpc request",
				zap.String("method", method), zap.Reflect("req", req), zap.Reflect("reply", reply),
				zap.Duration("cost", time.Since(begin)), zap.Error(err),
			)
		}()
		err = invoker(ctx, method, req, reply, cc, opts...)
		return err
	}
}

func withTimeout(timeout time.Duration) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func withTracer() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		commonCtx := ctm_context.GetContext(ctx)
		if commonCtx.Span == nil {
			return invoker(ctx, method, req, reply, cc, opts...)
		}

		span := commonCtx.Span.Tracer().StartSpan(fmt.Sprintf("grpc_client:%s", method), opentracing.ChildOf(commonCtx.Span.Context()))
		defer span.Finish()

		carrier := opentracing.TextMapCarrier{}
		if err := opentracing.GlobalTracer().Inject(commonCtx.Span.Context(), opentracing.HTTPHeaders, &carrier); err != nil {
			commonCtx.Logger.Warn("Extract span fail", zap.Error(err))
		}
		md := metadata.MD{}
		for key, value := range carrier {
			md.Set(key, value)
		}

		defer func() {
			span.SetTag("req", req)
			span.SetTag("reply", reply)
		}()

		ctx = metadata.NewOutgoingContext(ctx, md)
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
