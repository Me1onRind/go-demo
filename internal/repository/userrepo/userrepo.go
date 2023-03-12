package userrepo

import (
	"context"

	"github.com/Me1onRind/go-demo/internal/global/gconfig"
	"github.com/Me1onRind/go-demo/internal/infrastructure/client/mysql"
	"github.com/Me1onRind/go-demo/internal/infrastructure/logger"
	"github.com/Me1onRind/go-demo/internal/model/po/userpo"
	"gorm.io/gorm"
)

type UserRepo struct {
}

func NewUserRepo() *UserRepo {
	return &UserRepo{}
}

func (u *UserRepo) DBLabel() string {
	return gconfig.DynamicCfg.DefaultDB.GetLabel()
}

func (u *UserRepo) GetUser(ctx context.Context, opts ...mysql.Option) (*userpo.User, error) {
	logger.CtxInfof(ctx, "get user")
	return mysql.Take[userpo.User](ctx, u.DBLabel(), opts...)
}

func (u *UserRepo) CreateUser(ctx context.Context, user *userpo.User) (*userpo.User, error) {
	logger.CtxInfof(ctx, "create user")
	return mysql.Create(ctx, u.DBLabel(), user)
}

func WithUserId(userId uint64) mysql.Option {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("user_id=?", userId)
	}
}

func WithEmail(email string) mysql.Option {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("email=?", email)
	}
}
