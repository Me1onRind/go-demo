package foo_service

import (
	"github.com/Me1onRind/go-demo/internal/core/common"
	"github.com/Me1onRind/go-demo/internal/third_party/go_demo_service"
)

type FooService struct {
	GoDemoService *go_demo_service.GoDemoService
}

func NewFooService() *FooService {
	return &FooService{
		GoDemoService: go_demo_service.NewGoDemoService(),
	}
}

func (f *FooService) ProxyGreet(ctx *common.Context, name, msg string) (string, *common.Error) {
	return f.GoDemoService.Greet(ctx, name, msg)
}
