package store

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/Me1onRind/go-demo/internal/core/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DB       *gorm.DB
	ConfigDB *gorm.DB
)

func NewDBPoolFromDsn(dsn string) (*gorm.DB, error) {
	return doCreateDBPool(mysql.Open(dsn))
}

func NewDBPool(cfg *config.MysqlConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?timeout=%s&readTimeout=%s&writeTimeout=%s",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.Timeout, cfg.ReadTimeout, cfg.WriteTimeout)
	dail := mysql.Open(dsn)

	db, err := doCreateDBPool(dail)
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxIdleTime(time.Second * cfg.MaxIdleTime)
	sqlDB.SetConnMaxLifetime(time.Second * cfg.MaxLifetime)
	return db, nil
}

func NewDBPoolFromDB(db *sql.DB) (*gorm.DB, error) {
	return doCreateDBPool(mysql.New(mysql.Config{
		Conn: db,
	}))
}

func doCreateDBPool(dial gorm.Dialector) (*gorm.DB, error) {
	db, err := gorm.Open(dial, &gorm.Config{
		//Logger: logger.Default.LogMode(logger.Info),
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}

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
			return
		}
	}
}

func tracingEnd() func(*gorm.DB) {
	return func(db *gorm.DB) {
		if db.Statement.Context != nil {
			return
		}
	}
}
