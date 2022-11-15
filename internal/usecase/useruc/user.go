package useruc

import "context"

type UserUsecase struct {
}

func NewUserUsecase() *UserUsecase {
	return &UserUsecase{}
}

func (u *UserUsecase) QueryUserInfo(ctx context.Context, raw any) (any, error) {
	return map[string]string{
		"user_id": "mock",
	}, nil
}
