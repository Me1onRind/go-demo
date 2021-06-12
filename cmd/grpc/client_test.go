package main

import (
	"context"
	"testing"
	"time"

	"github.com/Me1onRind/go-demo/internal/core/common"
	"github.com/Me1onRind/go-demo/internal/core/middleware"
	"github.com/Me1onRind/go-demo/protobuf/pb"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/resolver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

var (
	c         pb.FooClient
	commonCtx = common.NewContext(context.Background())
)

func TestMain(m *testing.M) {
	cli, err := clientv3.NewFromURL("http://localhost:2379")
	if err != nil {
		panic(err)
	}
	builder, err := resolver.NewBuilder(cli)
	if err != nil {
		panic(err)
	}
	conn, err := grpc.Dial("etcd:///service/go-demo",
		grpc.WithResolvers(builder),
		grpc.WithBalancerName("round_robin"),
		grpc.WithInsecure(), grpc.WithTimeout(time.Second),
		grpc.WithChainUnaryInterceptor(grpc_middleware.ChainUnaryClient(
			middleware.ClientRetry(1),
			middleware.ClientLogger(),
		)),
	)
	if err != nil {
		panic(err)
	}
	c = pb.NewFooClient(conn)
	m.Run()
}

func Test_Greet(t *testing.T) {
	t.Skip()
	ctx, cancel := context.WithTimeout(commonCtx, time.Second*2)
	defer cancel()
	resp, err := c.Greet(ctx, &pb.GreetReq{
		MyName: "Bar",
		Msg:    "Hello, World",
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(resp.Msg)
}

func Test_ErrorResult(t *testing.T) {
	ctx, cancel := context.WithTimeout(commonCtx, time.Second*2)
	defer cancel()
	_, err := c.ErrorResult(ctx, &pb.Empty{})
	if err != nil {
		if s, ok := status.FromError(err); ok {
			details := s.Details()
			if len(details) > 0 {
				if e, ok := details[0].(*pb.Error); ok {
					t.Log(e.Code, e.Message)
				}
			}
		} else {
			t.Fatal(err)
		}
	}
}

func Test_PanicResult(t *testing.T) {
	ctx, cancel := context.WithTimeout(commonCtx, time.Second*2)
	defer cancel()
	_, err := c.PanicResult(ctx, &pb.Empty{})
	t.Log(err)
}
