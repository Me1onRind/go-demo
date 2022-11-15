package unexpectuc

import "context"

type UnexcpectUsecase struct {
}

func NewUnexpectUseCase() *UnexcpectUsecase {
	return &UnexcpectUsecase{}
}

func (u *UnexcpectUsecase) Panic(ctx context.Context, raw any) (any, error) {
	panic("panic api")
}
