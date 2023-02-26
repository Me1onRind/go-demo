package unexpectuc

import (
	"context"

	"github.com/Me1onRind/go-demo/protocol"
)

type UnexcpectUsecase struct {
}

func NewUnexpectUseCase() *UnexcpectUsecase {
	return &UnexcpectUsecase{}
}

func (u *UnexcpectUsecase) Panic(ctx context.Context, request *protocol.EmptyReq) (any, error) {
	panic("panic api")
}
