package foo_controller

import (
	"context"
	"fmt"
	"time"

	"github.com/Me1onRind/go-demo/internal/core/common"
	"github.com/Me1onRind/go-demo/internal/err_code"
	"github.com/Me1onRind/go-demo/internal/service/foo_service"
	"github.com/Me1onRind/go-demo/protobuf/pb"
	"github.com/Me1onRind/go-demo/protocol"
)

type FooController struct {
	FooService *foo_service.FooService
}

func NewFooController() *FooController {
	f := &FooController{
		FooService: foo_service.NewFooService(),
	}
	return f
}

func (f *FooController) Greet(ctx context.Context, in *pb.GreetReq) (*pb.GreetResp, error) {
	reply := fmt.Sprintf("Hello %s, I got your msg:%s", in.GetMyName(), in.GetMsg())
	out := &pb.GreetResp{}
	out.Msg = reply
	time.Sleep(time.Millisecond * 5)
	return out, nil
}

func (f *FooController) ErrorResult(ctx context.Context, in *pb.Empty) (*pb.Empty, error) {
	time.Sleep(time.Millisecond * 2)
	return nil, err_code.ServerInternalError.GrpcErr()
}

func (f *FooController) PanicResult(ctx context.Context, in *pb.Empty) (*pb.Empty, error) {
	time.Sleep(time.Millisecond * 2)
	panic("no implement")
}

func (f *FooController) ProxyGreet(ctx *common.Context, raw interface{}) (interface{}, *common.Error) {
	request := raw.(*protocol.GreetProxyRequest)
	reply, err := f.FooService.ProxyGreet(ctx, request.Name, request.Msg)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"reply": reply,
	}, nil
}
