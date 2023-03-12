package mysql

import (
	"context"

	"github.com/Me1onRind/go-demo/internal/global/gerror"
)

func Create[T any](ctx context.Context, dbLabel string, value *T, opts ...Option) (*T, error) {
	db := GetWriteDB(ctx, dbLabel, opts...).Create(value)
	if err := db.Error; err != nil {
		return nil, gerror.WriteDBError.Wrap(err)
	}
	return value, nil
}

func Find[T any](ctx context.Context, dbLabel string, opts ...Option) ([]*T, error) {
	result := []*T{}
	db := GetReadDB(ctx, dbLabel, opts...).Find(&result)
	if err := db.Error; err != nil {
		return nil, gerror.ReadDBError.Wrap(err)
	}
	return result, nil
}

func Take[T any](ctx context.Context, dbLabel string, opts ...Option) (*T, error) {
	var value T
	db := GetReadDB(ctx, dbLabel, opts...).Take(&value)
	if err := db.Error; err != nil {
		return nil, gerror.ReadDBError.Wrap(err)
	}
	return &value, nil
}
