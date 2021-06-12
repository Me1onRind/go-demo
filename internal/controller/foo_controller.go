package controller

import (
	"context"
	"fmt"
	"time"

	"github.com/Me1onRind/go-demo/internal/err_code"
	"github.com/Me1onRind/go-demo/protobuf/pb"
)

type FooController struct {
}

func NewFooController() *FooController {
	f := &FooController{}
	return f
}

func (f *FooController) Greet(ctx context.Context, in *pb.GreetReq) (*pb.GreetResp, error) {
	reply := fmt.Sprintf("Hello %s, I got your msg:%s", in.GetMyName(), in.GetMsg())
	out := &pb.GreetResp{}
	out.Msg = reply
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
