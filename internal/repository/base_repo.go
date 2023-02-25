package respository

import (
	"context"

	"github.com/Me1onRind/go-demo/internal/global/gerror"
	"github.com/Me1onRind/go-demo/internal/infrastructure/client/mysql"
	"github.com/Me1onRind/go-demo/internal/infrastructure/logger"
	"gorm.io/gorm"
)

type ModelIfce interface {
	TableName() string
	DBLabel() string
	GetId() uint64
}

type Option func(*gorm.DB) *gorm.DB

type BaseRepo[T ModelIfce] struct {
}

func NewBaseRepo[T ModelIfce]() *BaseRepo[T] {
	return &BaseRepo[T]{}
}

func (r *BaseRepo[T]) Create(ctx context.Context, value T) (T, error) {
	if err := mysql.GetWriteDB(ctx, value.DBLabel()).Create(value).Error; err != nil {
		logger.CtxErrorf(ctx, "Insert DB fail, db:[%s], table:[%s], value:[%+v], case:[%s]",
			value.DBLabel(), value.TableName(), value, err)
		var t T
		return t, gerror.WriteDBError.Wrap(err)
	}
	logger.CtxInfof(ctx, "Insert DB success, db:[%s], table:[%s], id:[%d]", value.DBLabel(), value.TableName(), value.GetId())
	return value, nil
}

func (r *BaseRepo[T]) Find(ctx context.Context, opts ...Option) ([]T, error) {
	result := []T{}
	var model T
	db := mysql.GetReadDB(ctx, model.DBLabel())
	for _, opt := range opts {
		opt(db)
	}
	if err := db.Find(&result).Error; err != nil {
		logger.CtxErrorf(ctx, "Find DB fail, db:[%s], table:[%s], case:[%s]",
			model.DBLabel(), model.TableName(), err)
		return nil, gerror.ReadDBError.Wrap(err)
	}
	return result, nil
}

func (r *BaseRepo[T]) Take(ctx context.Context, opts ...Option) (T, error) {
	var value T
	db := mysql.GetReadDB(ctx, value.DBLabel())
	for _, opt := range opts {
		db = opt(db)
	}
	if err := db.Take(&value).Error; err != nil {
		logger.CtxErrorf(ctx, "Take DB fail, db:[%s], table:[%s], case:[%s]",
			value.DBLabel(), value.TableName(), err)
		return value, gerror.ReadDBError.Wrap(err)
	}
	return value, nil
}
