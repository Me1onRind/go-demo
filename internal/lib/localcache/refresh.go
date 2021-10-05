package localcache

import (
	"github.com/Me1onRind/go-demo/internal/core/common"
	"github.com/Me1onRind/go-demo/internal/core/store"
	"github.com/Me1onRind/go-demo/internal/err_code"
	"go.uber.org/zap"
)

func Refresh(ctx *common.Context, localCacheKey string) *common.Error {
	version, err := store.RedisPool.Incr(ctx, versionKey(localCacheKey)).Uint64()
	ctx.Logger.Info("Refresh localcache version", zap.String("localCacheKey", localCacheKey), zap.Uint64("newVersion", version))
	if err != nil {
		ctx.Logger.Error("Refresh localcache version failed", zap.String("localCacheKey", localCacheKey))
		return err_code.WriteRedisError.WithErr(err)
	}
	return nil
}
