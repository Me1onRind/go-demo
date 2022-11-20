package mysql

import (
	"context"
	"fmt"

	"github.com/Me1onRind/go-demo/internal/infrastructure/logger"
	"github.com/Me1onRind/go-demo/internal/model/configmd"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

var (
	dbs = map[string]*gorm.DB{}
)

type TransactionHandler func(context.Context) error
type ctxTransactionKey string

func NewMysqlClusterClient(cfg *configmd.DBCluster) (*gorm.DB, error) {
	source := mysql.Open(cfg.Master.DSN)

	replicas := []gorm.Dialector{}
	for _, slave := range cfg.Slaves {
		replicas = append(replicas, mysql.Open(slave.DSN))
	}

	db, err := newMysqlClusterClientByDialector(source, replicas)
	if err != nil {
		return nil, err
	}

	if err := setDBLabel(db, cfg.Label); err != nil {
		return nil, err
	}

	return db, nil
}

func GetDB(ctx context.Context, dbLabel string) *gorm.DB {
	labelKey := ctxTransactionKey(dbLabel)
	tx := ctx.Value(labelKey)
	if tx != nil {
		return tx.(*gorm.DB)
	}

	return dbs[dbLabel]
}

func Transaction(ctx context.Context, dbLabel string, f TransactionHandler) (err error) {
	labelKey := ctxTransactionKey(dbLabel)
	t := ctx.Value(labelKey)

	// allow nested transaction
	if t != nil {
		return f(ctx)
	}

	db := dbs[dbLabel]
	if db == nil {
		return fmt.Errorf("DB label:[%s] is not found", dbLabel)
	}

	isPanic := true
	tx := db.Begin()
	ctx = context.WithValue(ctx, labelKey, tx)
	defer func() {
		if isPanic || err != nil {
			if rbErr := tx.Rollback().Error; rbErr != nil {
				// only log, keep origin error
				logger.CtxErrorf(ctx, "Transaction rollback err:[%s]", rbErr)
			}
			return
		}

		if cmErr := tx.Commit().Error; cmErr != nil {
			logger.CtxErrorf(ctx, "Transaction commit err:[%s]", cmErr)
			err = cmErr
		}
	}()

	err = f(ctx)
	isPanic = false
	return err
}

func setDBLabel(db *gorm.DB, label string) error {
	if len(label) == 0 {
		return fmt.Errorf("DB label is empty")
	}

	if _, ok := dbs[label]; ok {
		return fmt.Errorf("DB label:[%s] is existed", label)
	}
	dbs[label] = db
	return nil
}

func newMysqlClusterClientByDialector(source gorm.Dialector, replicas []gorm.Dialector) (*gorm.DB, error) {
	db, err := gorm.Open(source, &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		return nil, err
	}

	if len(replicas) > 0 {
		if err := db.Use(dbresolver.Register(dbresolver.Config{
			Replicas: replicas,
			Policy:   dbresolver.RandomPolicy{},
		})); err != nil {
			return nil, err
		}
	}
	return db, nil
}
