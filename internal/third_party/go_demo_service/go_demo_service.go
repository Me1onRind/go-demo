package go_demo_service

import (
	"github.com/Me1onRind/go-demo/err_code"
	"github.com/Me1onRind/go-demo/infrastructure/ctm_context"
)

type GoDemoService struct {
}

func NewGoDemoService() *GoDemoService {
	return &GoDemoService{}
}

func (g *GoDemoService) Greet(ctx *ctm_context.Context, name, msg string) (string, *err_code.Error) {

	return "ok", nil
}
