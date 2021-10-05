package initialize

import (
	"github.com/Me1onRind/go-demo/global/store"
	"github.com/Me1onRind/go-demo/global/store/db_label"
	"github.com/Me1onRind/go-demo/internal/core/config"
	"github.com/Me1onRind/go-demo/internal/lib/client/mysql_client"
	"github.com/bluele/gcache"
)

func InitDB() error {
	var err error
	dbs := config.RemoteConfig.DBs
	store.DBs[db_label.DB], err = mysql_client.NewDBClient(&dbs.DB)
	if err != nil {
		return err
	}
	store.DBs[db_label.ConfigDB], err = mysql_client.NewDBClient(&dbs.ConfigDB)
	if err != nil {
		return err
	}
	return nil
}

func CloseDB() error {
	for _, v := range store.DBs {
		db, err := v.DB()
		if err != nil {
			return err
		}
		if err := db.Close(); err != nil {
			return err
		}
	}
	return nil
}

func InitLocalCache() error {
	store.ConfigCache = gcache.New(100000).Build()
	return nil
}
