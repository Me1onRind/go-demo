package userus

import "github.com/Me1onRind/go-demo/internal/infrastructure/base"

type UserUsecase struct {
}

func NewUserUsecase() *UserUsecase {
	return &UserUsecase{}
}

func (u *UserUsecase) QueryUserInfo(ctx *base.Context, raw interface{}) (interface{}, error) {
	return map[string]string{
		"user_id": "mock",
	}, nil
}
