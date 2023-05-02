package pinguc

import (
	"context"

	"github.com/Me1onRind/go-demo/internal/infrastructure/logger"
	"github.com/Me1onRind/go-demo/protocol/pingproto"
)

type PingUsecase struct {
}

func NewPingUsecase() *PingUsecase {
	return &PingUsecase{}
}

func (p *PingUsecase) Ping(ctx context.Context, proto *pingproto.Ping) error {
	logger.CtxInfof(ctx, proto.Value)
	return nil
}
