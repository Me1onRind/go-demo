package localcache

import (
	"github.com/Me1onRind/go-demo/internal/lib/ctm_context"
	"github.com/Me1onRind/go-demo/internal/lib/err_code"
	"github.com/bluele/gcache"
)

type Loader interface {
	GetLocalcacheKey() string
	SetLocalcache(cache gcache.Cache)
	LocalcacheVersion(ctx *ctm_context.Context) (uint64, *err_code.Error)
	RefreshLocalcacheVersion(ctx *ctm_context.Context) *err_code.Error
	LoadLocalcacheData(ctx *ctm_context.Context) (interface{}, *err_code.Error)
}
