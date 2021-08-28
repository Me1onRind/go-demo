package store

import (
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewConnectPool(dns string) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(dns), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetConnMaxIdleTime(time.Second * 10)
	sqlDB.SetConnMaxLifetime(time.Minute * 3)

	registerPlugin(db)

	return db, err
}

func registerPlugin(db *gorm.DB) {
	_ = db.Callback().Query().Before("*").Register("my_plugin:tracing_start", tracingStart())
	_ = db.Callback().Query().After("*").Register("my_plugin:tracing_end", tracingEnd())
	_ = db.Callback().Update().Before("*").Register("my_plugin:set_mtime", setMTime)
	_ = db.Callback().Create().Before("*").Register("my_plugin:set_ctime_and_mtime", setCTimeAndMTime)
}

func setCTimeAndMTime(db *gorm.DB) {
	now := time.Now().Unix()
	db.Statement.SetColumn("ctime", now)
	db.Statement.SetColumn("mtime", now)
}

func setMTime(db *gorm.DB) {
	db.Statement.SetColumn("mtime", time.Now().Unix())
}

func tracingStart() func(*gorm.DB) {
	return func(db *gorm.DB) {
		if db.Statement.Context != nil {
		}
	}
}

func tracingEnd() func(*gorm.DB) {
	return func(db *gorm.DB) {
		if db.Statement.Context != nil {
		}
	}
}
