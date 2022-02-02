package grpc_client

import (
	"context"
	"time"

	"github.com/Me1onRind/go-demo/err_code"
	"github.com/Me1onRind/go-demo/protobuf/pb"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	gresolver "google.golang.org/grpc/resolver"
	"google.golang.org/grpc/status"
)

var (
	FooClient pb.FooClient
	fooCoon   *grpc.ClientConn
)

func InitGoDemoClient(target string, resolver gresolver.Builder) error {
	var err error
	conn, err := getGrpcConn(target, resolver)
	if err != nil {
		return err
	}

	FooClient = pb.NewFooClient(conn)
	return nil
}

func CloseGoDemoClient() error {
	return fooCoon.Close()
}

func ResolveGoDemoGrpcError(err error) *err_code.Error {
	if s, ok := status.FromError(err); ok {
		details := s.Details()
		if len(details) > 0 {
			if e, ok := details[0].(*pb.Error); ok {
				return err_code.GoDemoCommonFailedError.Withf("remote code:%d, message:%s", e.Code, e.Message)
			}
		}
	}
	return err_code.GRPCCallFailedError.WithErr(err)
}

func getGrpcConn(target string, resolver gresolver.Builder) (*grpc.ClientConn, error) {
	//func getGrpcConn(serviceName string) (*grpc.ClientConn, error) {
	//resolver, err := resolver.NewBuilder(etcd_client.EtcdClient)
	//if err != nil {
	//return nil, err
	//}

	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*2)
	defer cancel()

	policy := `{"loadBalancingPolicy":"round_robin"}`
	conn, err := grpc.DialContext(ctx, target,
		grpc.WithResolvers(resolver),
		grpc.WithDefaultServiceConfig(policy),
		grpc.WithInsecure(),
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
