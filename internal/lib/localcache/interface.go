package localcache

import (
	"github.com/Me1onRind/go-demo/internal/lib/ctm_context"
	"github.com/Me1onRind/go-demo/internal/lib/err_code"
	"github.com/bluele/gcache"
)

type Loader interface {
	LoadLocalCacheData(ctx *ctm_context.Context) (interface{}, *err_code.Error)
	LocalCacheKey() string
	LocalCacheInstance() gcache.Cache
}
