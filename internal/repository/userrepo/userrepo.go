package userrepo

import (
	"context"

	"github.com/Me1onRind/go-demo/internal/model/po/userpo"
	respository "github.com/Me1onRind/go-demo/internal/repository"
	"gorm.io/gorm"
)

type UserRepo struct {
	baseRepo *respository.BaseRepo[*userpo.User]
}

func NewUserRepo() *UserRepo {
	return &UserRepo{
		baseRepo: respository.NewBaseRepo[*userpo.User](),
	}
}

func (u *UserRepo) GetUser(ctx context.Context, opts ...respository.Option) (*userpo.User, error) {
	return u.baseRepo.Take(ctx, opts...)
}

func (u *UserRepo) CreateUser(ctx context.Context, user *userpo.User) (*userpo.User, error) {
	return u.baseRepo.Create(ctx, user)
}

func WithUserId(userId uint64) respository.Option {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("user_id=?", userId)
	}
}

func WithEmail(email string) respository.Option {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("email=?", email)
	}
}
