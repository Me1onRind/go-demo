package useruc

import (
	"context"
	"github.com/Me1onRind/go-demo/protocol/userproto"
)

type UserUsecase struct {
}

func NewUserUsecase() *UserUsecase {
	return &UserUsecase{}
}

func (u *UserUsecase) GetUserDetail(ctx context.Context, raw any) (any, error) {
	request := raw.(*userproto.GetUserDetailReq)
	return map[string]any{
		"user_id": request.UserId,
	}, nil
}
