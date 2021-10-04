package initialize

import (
	"time"

	"github.com/Me1onRind/go-demo/internal/core/config"
	"github.com/Me1onRind/go-demo/internal/core/store"
	"github.com/bluele/gcache"
)

func InitDB() error {
	var err error
	dbs := config.RemoteConfig.DBs
	store.DB, err = store.NewDBPool(&dbs.DB)
	if err != nil {
		return err
	}
	store.ConfigDB, err = store.NewDBPool(&dbs.ConfigDB)
	if err != nil {
		return err
	}
	return nil
}

func InitLocalCache() error {
	store.ConfigCache = gcache.New(100000).Expiration(time.Minute * 30).Build()
	return nil
}
