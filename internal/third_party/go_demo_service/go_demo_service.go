package go_demo_service

import (
	"github.com/Me1onRind/go-demo/internal/core/client/grpc_client"
	"github.com/Me1onRind/go-demo/internal/core/common"
	"github.com/Me1onRind/go-demo/protobuf/pb"
)

type GoDemoService struct {
}

func NewGoDemoService() *GoDemoService {
	return &GoDemoService{}
}

func (g *GoDemoService) Greet(ctx *common.Context, name, msg string) (string, *common.Error) {
	resp, err := grpc_client.FooClient.Greet(ctx, &pb.GreetReq{
		MyName: name,
		Msg:    msg,
	})
	if err != nil {
		return "", grpc_client.ResolveGoDemoGrpcError(err)
	}

	return resp.Msg, nil
}
