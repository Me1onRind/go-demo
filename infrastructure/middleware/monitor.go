package middleware

import (
	"context"

	"github.com/Me1onRind/go-demo/global/prometheus_vec"
	"google.golang.org/grpc"
)

func GrpcMonitor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		prometheus_vec.ReqTotalCounterVec.WithLabelValues(info.FullMethod).Inc()
		return handler(ctx, req)
	}
}
