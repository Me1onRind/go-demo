package localcache

import (
	"github.com/Me1onRind/go-demo/global/store"
	"github.com/Me1onRind/go-demo/internal/lib/ctm_context"
	"github.com/Me1onRind/go-demo/internal/lib/err_code"
	"go.uber.org/zap"
)

func Refresh(ctx *ctm_context.Context, localCacheKey string) *err_code.Error {
	version, err := store.RedisClient.Incr(ctx, versionKey(localCacheKey)).Uint64()
	ctx.Logger.Info("Refresh localcache version", zap.String("localCacheKey", localCacheKey), zap.Uint64("newVersion", version))
	if err != nil {
		ctx.Logger.Error("Refresh localcache version failed", zap.String("localCacheKey", localCacheKey))
		return err_code.WriteRedisError.WithErr(err)
	}
	return nil
}
