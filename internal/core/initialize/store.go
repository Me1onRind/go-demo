package initialize

import (
	"time"

	"github.com/Me1onRind/go-demo/internal/core/store"
	"github.com/bluele/gcache"
)

func InitDB() error {
	var err error
	store.DB, err = store.NewDBConnectPool("root:guapi123@tcp(localhost:3306)/go-frame")
	if err != nil {
		return err
	}
	store.ConfigDB, err = store.NewDBConnectPool("root:guapi123@tcp(localhost:3306)/config")
	return err
}

func InitLocalCache() error {
	store.ConfigCache = gcache.New(100000).Expiration(time.Minute * 30).Build()
	return nil
}
