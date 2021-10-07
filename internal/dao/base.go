package dao

import (
	"errors"
	"fmt"

	"github.com/Me1onRind/go-demo/global/store"
	"github.com/Me1onRind/go-demo/internal/core/logger"
	"github.com/Me1onRind/go-demo/internal/lib/ctm_context"
	"github.com/Me1onRind/go-demo/internal/lib/err_code"
	"github.com/bluele/gcache"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

type LocalcacheDao struct {
	Localcache    gcache.Cache
	LocalcacheKey string
}

func (l *LocalcacheDao) SetLocalcache(cache gcache.Cache) {
	l.Localcache = cache
}

func (l *LocalcacheDao) GetLocalcacheKey() string {
	return l.LocalcacheKey
}

func (l *LocalcacheDao) LocalcacheVersion(ctx *ctm_context.Context) (uint64, *err_code.Error) {
	version, err := store.RedisClient.Get(ctx, l.redisVersionKey()).Uint64()
	if err != nil && !errors.Is(err, redis.Nil) {
		logger.StdoutLogger.Error("Get localcache version failed", zap.String("localcacheKey", l.LocalcacheKey), zap.Error(err))
		return 0, err_code.ReadRedisError.WithErr(err)
	}

	return version, nil
}

func (l *LocalcacheDao) RefreshLocalcacheVersion(ctx *ctm_context.Context) *err_code.Error {
	version, err := store.RedisClient.Incr(ctx, l.redisVersionKey()).Uint64()
	ctx.Logger.Info("Refresh localcache version", zap.String("localcacheKey", l.LocalcacheKey), zap.Uint64("newVersion", version))
	if err != nil {
		ctx.Logger.Error("Refresh localcache version failed", zap.String("localCacheKey", l.LocalcacheKey))
		return err_code.WriteRedisError.WithErr(err)
	}
	return nil
}

func (l *LocalcacheDao) redisVersionKey() string {
	return fmt.Sprintf("localcache:version:%s", l.LocalcacheKey)
}
