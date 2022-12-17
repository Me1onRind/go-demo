package mysql

import (
	"context"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Me1onRind/go-demo/internal/infrastructure/logger"
	"github.com/Me1onRind/go-demo/internal/model/configmd"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

var (
	dbs *gormSet = newGormSet()
)

type TransactionHandler func(context.Context) error
type ctxTransactionKey string

func RegisterMysqlCluster(cfg *configmd.DBCluster) error {
	master, err := newGormDBByDialector(mysql.Open(cfg.Master.DSN), &dbDesc{
		role:    "master",
		dbLabel: cfg.Label,
	})
	if err != nil {
		return err
	}

	replicas := []*gorm.DB{}
	for _, slave := range cfg.Slaves {
		replica, err := newGormDBByDialector(mysql.Open(slave.DSN), &dbDesc{
			role:    "slave",
			dbLabel: cfg.Label,
		})
		if err != nil {
			return err
		}
		replicas = append(replicas, replica)
	}

	dbs.register(cfg.Label, newGormDBInfo(master, replicas))
	return nil
}

func GetWriteDB(ctx context.Context, label string) *gorm.DB {
	db, _ := getTxOrWriteDB(ctx, label)
	return db
}

func GetReadDB(ctx context.Context, label string) *gorm.DB {
	tx := getTx(ctx, label)
	if tx != nil {
		return tx
	}
	return dbs.getGormDBInfo(label).getReadDB().WithContext(ctx)
}

func Transaction(ctx context.Context, label string, f TransactionHandler) (err error) {
	tx, isTransaction := getTxOrWriteDB(ctx, label)
	if isTransaction {
		return f(ctx)
	}
	tx = tx.Begin()
	ctx = context.WithValue(ctx, ctxTransactionKey(label), tx)

	isPanic := true
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

func NewMysqlMock(label string) sqlmock.Sqlmock {
	sqlDb, mock, _ := sqlmock.New()
	mock.ExpectQuery("SELECT VERSION").WillReturnRows(sqlmock.NewRows([]string{"VERSION"}).AddRow("5.7.32"))
	masterResource := mysql.New(mysql.Config{Conn: sqlDb})
	db, _ := newGormDBByDialector(masterResource, &dbDesc{role: "mock", dbLabel: label})
	dbs.register(label, newGormDBInfo(db, nil))
	return mock
}

func newGormDBByDialector(source gorm.Dialector, desc *dbDesc) (*gorm.DB, error) {
	db, err := gorm.Open(source, &gorm.Config{
		SkipDefaultTransaction: true,
		Logger:                 gormLogger.Default.LogMode(gormLogger.Silent),
	})
	if err != nil {
		return nil, err
	}
	registerPlugin(db, desc)
	return db, nil
}

func getTxOrWriteDB(ctx context.Context, label string) (*gorm.DB, bool) {
	tx := getTx(ctx, label)
	if tx != nil {
		return tx, true
	}
	return dbs.getGormDBInfo(label).getWriteDB().WithContext(ctx), false
}

func getTx(ctx context.Context, label string) *gorm.DB {
	tx := ctx.Value(ctxTransactionKey(label))
	if tx != nil {
		return tx.(*gorm.DB)
	}
	return nil
}

//func GetWriteDB(ctx context.Context, dbLabel string) *gorm.DB {
//db, _ := getDB(ctx, dbLabel)
//return db
//}

//func getDB(ctx context.Context, dbLabel string) (*gorm.DB, bool) {
//tx := ctx.Value(ctxTransactionKey(dbLabel))
//if tx != nil {
//return tx.(*gorm.DB), true
//}

//db := dbs[dbLabel]
//if db == nil {
//panic(fmt.Sprintf("DB label:[%s] not exist", dbLabel))
//}
//return db.WithContext(ctx), false
//}

//func setDBLabel(db *gorm.DB, label string, checkDuplicate bool) error {
//if len(label) == 0 {
//return fmt.Errorf("DB label is empty")
//}

//if _, ok := dbs[label]; checkDuplicate && ok {
//return fmt.Errorf("DB label:[%s] is existed", label)
//}
//dbs[label] = db
//return nil
//}
