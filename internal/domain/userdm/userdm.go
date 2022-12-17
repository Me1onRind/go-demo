package userdm

import (
	"context"
	"errors"

	"github.com/Me1onRind/go-demo/internal/domain/iddm"
	"github.com/Me1onRind/go-demo/internal/global/gerror"
	"github.com/Me1onRind/go-demo/internal/model/po/idpo"
	"github.com/Me1onRind/go-demo/internal/model/po/userpo"
	"github.com/Me1onRind/go-demo/internal/repository/userrepo"
	"gorm.io/gorm"
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

	_, err := u.UserRepo.GetUserByEmail(ctx, user.Email)
	if err == nil {
		return nil, gerror.RecordExistedError.Withf("email:[%s]", user.Email)
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	maxTry := 5
	userId, err := u.IdDomain.GetId(ctx, idpo.UserIdType, maxTry)
	if err != nil {
		return nil, err
	}

	user.UserId = userId
	return u.UserRepo.CreateUser(ctx, user)
}
