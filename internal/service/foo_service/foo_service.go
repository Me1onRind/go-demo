package foo_service

import (
	//"github.com/Me1onRind/go-demo/internal/lib/ctm_context"
	"github.com/Me1onRind/go-demo/err_code"
	"github.com/Me1onRind/go-demo/infrastructure/ctm_context"
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

func (f *FooService) ProxyGreet(ctx *ctm_context.Context, name, msg string) (string, *err_code.Error) {
	return f.GoDemoService.Greet(ctx, name, msg)
}
