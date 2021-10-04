package main

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/Me1onRind/go-demo/internal/core/client/grpc_client"
	"github.com/Me1onRind/go-demo/internal/core/common"
	"github.com/Me1onRind/go-demo/internal/core/initialize"
	"github.com/Me1onRind/go-demo/protobuf/pb"
	"google.golang.org/grpc/status"
)

var (
	commonCtx *common.Context
)

func TestMain(m *testing.M) {
	_ = initialize.InitLogger()
	commonCtx = common.NewContext(context.Background())
	if err := initialize.InitGrpcClients(); err != nil {
		return
	}
	os.Exit(m.Run())
}

func Test_Greet(t *testing.T) {
	for i := 0; i < 100; i++ {
		ctx, cancel := context.WithTimeout(commonCtx, time.Second*2)
		defer cancel()
		resp, err := grpc_client.FooClient.Greet(ctx, &pb.GreetReq{
			MyName: "Bar",
			Msg:    "Hello, World",
		})
		if err != nil {
			t.Fatal(err)
		}
		t.Log(resp.Msg)
	}
}

func Test_ErrorResult(t *testing.T) {
	ctx, cancel := context.WithTimeout(commonCtx, time.Second*2)
	defer cancel()
	_, err := grpc_client.FooClient.ErrorResult(ctx, &pb.Empty{})
	if err != nil {
		if s, ok := status.FromError(err); ok {
			details := s.Details()
			if len(details) > 0 {
				if e, ok := details[0].(*pb.Error); ok {
					t.Logf("pb return code:%d, msg:%s\n", e.Code, e.Message)
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
	_, err := grpc_client.FooClient.PanicResult(ctx, &pb.Empty{})
	t.Log(err)
}
