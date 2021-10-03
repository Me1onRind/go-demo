package localcache

import (
	"github.com/Me1onRind/go-demo/internal/core/common"
)

type Loader interface {
	LoadLocalCacheData() *common.Error
}
