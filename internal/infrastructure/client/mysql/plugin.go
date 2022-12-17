package mysql

import (
	"time"

	"github.com/Me1onRind/go-demo/internal/infrastructure/logger"
	"gorm.io/gorm"
)

const (
	startTimeKey = "start_time"
)

type dbDesc struct {
	dbLabel string
	role    string
}

func registerPlugin(db *gorm.DB, desc *dbDesc) {
	_ = db.Callback().Create().Before("*").Register("my_plugin:create_before_monitor", beforeMonitor(desc))
	_ = db.Callback().Delete().Before("*").Register("my_plugin:delete_before_monitor", beforeMonitor(desc))
	_ = db.Callback().Update().Before("*").Register("my_plugin:update_before_monitor", beforeMonitor(desc))
	_ = db.Callback().Query().Before("*").Register("my_plugin:query_before_monitor", beforeMonitor(desc))

	_ = db.Callback().Create().After("*").Register("my_plugin:create_after_monitor", afterMonitor(desc))
	_ = db.Callback().Delete().After("*").Register("my_plugin:delete_after_monitor", afterMonitor(desc))
	_ = db.Callback().Update().After("*").Register("my_plugin:update_after_monitor", afterMonitor(desc))
	_ = db.Callback().Query().After("*").Register("my_plugin:query_after_monitor", afterMonitor(desc))
}

func beforeMonitor(desc *dbDesc) func(*gorm.DB) {
	return func(db *gorm.DB) {
		startTime := time.Now()
		db.Set(startTimeKey, startTime)
	}
}

func afterMonitor(desc *dbDesc) func(*gorm.DB) {
	return func(db *gorm.DB) {
		statement := db.Statement
		ctx := statement.Context

		st, _ := db.Get(startTimeKey)
		startTime := st.(time.Time)
		duration := time.Since(startTime)

		//m := statement.Dest.(*mysql.TestModel)

		logger.CtxInfof(ctx, "db:[%s],role:[%s],sql:[%s],args:[%v],cost:[%s]",
			desc.dbLabel, desc.role, statement.SQL.String(), statement.Vars, duration)
	}
}
