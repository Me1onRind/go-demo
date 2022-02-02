package mysql_client

import (
	"database/sql"
	"time"

	"github.com/Me1onRind/go-demo/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
)

func NewDBClient(cfg *config.MysqlResources) (*gorm.DB, error) {
	masterCfg := &cfg.Master
	dail := mysql.Open(masterCfg.DSN(cfg.DBName))
	db, err := doCreateDBClient(dail)
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	masterPoolCfg := &cfg.MasterPool
	sqlDB.SetMaxIdleConns(masterPoolCfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(masterPoolCfg.MaxOpenConns)
	sqlDB.SetConnMaxIdleTime(time.Second * masterPoolCfg.MaxIdleTime)
	sqlDB.SetConnMaxLifetime(time.Second * masterPoolCfg.MaxLifetime)

	slaveDialector := []gorm.Dialector{}
	for _, v := range cfg.Slaves {
		slaveDialector = append(slaveDialector, mysql.Open(v.DSN(cfg.DBName)))
	}
	if len(slaveDialector) > 0 {
		slavePoolCfg := cfg.SlavePool
		err := db.Use(
			dbresolver.Register(dbresolver.Config{
				Replicas: slaveDialector,
				Policy:   dbresolver.RandomPolicy{},
			}).SetMaxIdleConns(slavePoolCfg.MaxIdleConns).
				SetMaxOpenConns(slavePoolCfg.MaxOpenConns).
				SetConnMaxIdleTime(slavePoolCfg.MaxIdleTime * time.Second).
				SetConnMaxLifetime(slavePoolCfg.MaxLifetime * time.Second),
		)
		if err != nil {
			return nil, err
		}
	}

	return db, nil
}

func NewDBClientFromDB(db *sql.DB) (*gorm.DB, error) {
	return doCreateDBClient(mysql.New(mysql.Config{
		Conn: db,
	}))
}

func doCreateDBClient(dial gorm.Dialector) (*gorm.DB, error) {
	db, err := gorm.Open(dial, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		//Logger:                 logger.Default.LogMode(logger.Silent),
		SkipDefaultTransaction: true,
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
