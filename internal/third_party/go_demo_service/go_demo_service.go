package go_demo_service

import (
	"github.com/Me1onRind/go-demo/infrastructure/ctm_context"
	"github.com/Me1onRind/go-demo/internal/lib/err_code"
)

type GoDemoService struct {
}

func NewGoDemoService() *GoDemoService {
	return &GoDemoService{}
}

func (g *GoDemoService) Greet(ctx *ctm_context.Context, name, msg string) (string, *err_code.Error) {

	return "ok", nil
}
