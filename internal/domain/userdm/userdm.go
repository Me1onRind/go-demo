package userdm

import (
	"context"

	"github.com/Me1onRind/go-demo/internal/domain/iddm"
	"github.com/Me1onRind/go-demo/internal/model/po/idpo"
	"github.com/Me1onRind/go-demo/internal/model/po/userpo"
	"github.com/Me1onRind/go-demo/internal/repository/userrepo"
)

type UserDomain struct {
	UserRepo *userrepo.UserRepo
	IdDomain iddm.IdDomain
}

func NewUserDomain() *UserDomain {
	return &UserDomain{
		IdDomain: iddm.NewIdDomain(),
		UserRepo: userrepo.NewUserRepo(),
	}
}

func (u *UserDomain) CreateUser(ctx context.Context, user *userpo.User) (*userpo.User, error) {
	maxTry := 5
	userId, err := u.IdDomain.GetId(ctx, idpo.UserIdType, maxTry)
	if err != nil {
		return nil, err
	}

	user.UserId = userId
	return u.UserRepo.CreateUser(ctx, user)
}
