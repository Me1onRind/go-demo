package grpc_client

import (
	"time"

	"github.com/Me1onRind/go-demo/internal/core/common"
	"github.com/Me1onRind/go-demo/internal/core/register"
	"github.com/Me1onRind/go-demo/internal/err_code"
	"github.com/Me1onRind/go-demo/protobuf/pb"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	FooClient pb.FooClient
)

func InitGoDemoClient() error {
	conn, err := getGrpcConn("go-demo")
	if err != nil {
		return err
	}

	FooClient = pb.NewFooClient(conn)
	return nil
}

func ResolveGoDemoGrpcError(err error) *common.Error {
	if s, ok := status.FromError(err); ok {
		details := s.Details()
		if len(details) > 0 {
			if e, ok := details[0].(*pb.Error); ok {
				if result := common.GetErrorByCode(e.Code); result != nil {
					return result
				}
				return err_code.GoDemoCommonFailedError.Withf("remote code:%d, message:%s", e.Code, e.Message)
			}
		}
	}
	return err_code.GRPCCallFailedError.WithErr(err)
}

func getGrpcConn(serviceName string) (*grpc.ClientConn, error) {
	resolver, err := register.GrpcResolvers()
	if err != nil {
		return nil, err
	}

	conn, err := grpc.Dial(register.DialTarget(serviceName),
		grpc.WithResolvers(resolver),
		grpc.WithBalancerName("round_robin"),
		grpc.WithInsecure(), grpc.WithTimeout(time.Second),
		grpc.WithChainUnaryInterceptor(grpc_middleware.ChainUnaryClient(
			withRetry(1, retryOnlyByCode(codes.Unavailable)),
			withTimeout(time.Second*1),
			withLogger(),
			withTracer(),
		)),
	)

	if err != nil {
		return nil, err
	}

	return conn, nil
}
