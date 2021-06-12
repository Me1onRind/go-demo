package middleware

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type RetryFilter func(method string, req, reply interface{}, err error) bool

func ClientRetry(maxRetry uint32, needRetry RetryFilter) grpc.UnaryClientInterceptor {
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

func RetryOnlyByCode(codes ...codes.Code) RetryFilter {
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
