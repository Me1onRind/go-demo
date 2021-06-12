package middleware

import (
	"context"
	"google.golang.org/grpc"
)

func ClientRetry(maxRetry uint32) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		var err error
		for i := uint32(0); i < maxRetry+1; i++ {
			err = invoker(ctx, method, req, reply, cc, opts...)
			if err == nil {
				return nil
			}
		}
		return err
	}
}
