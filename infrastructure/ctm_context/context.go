package ctm_context

import (
	"context"
	"fmt"

	"github.com/Me1onRind/go-demo/err_code"
	"github.com/Me1onRind/go-demo/global/client_singleton"
	"github.com/Me1onRind/go-demo/infrastructure/db_label"
	"gorm.io/gorm"
)

const (
	ctmCtxKey = "ctm"
)

type Context struct {
	context.Context
	txs map[db_label.Label]*gorm.DB
}

func NewContext(ctx context.Context) *Context {
	c := &Context{
		Context: ctx,
		txs:     map[db_label.Label]*gorm.DB{},
	}
	return c
}

func (c *Context) GetDB(label db_label.Label) *gorm.DB {
	tx := c.txs[label]
	if tx != nil {
		return tx
	}

	db := client_singleton.DBs[label]
	if db == nil {
		panic(fmt.Sprintf("Can't get DB, DB label:%s", label))
	}
	return db
}

func (c *Context) Transaction(label db_label.Label, fc func() *err_code.Error) (err *err_code.Error) {
	tx := c.txs[label]
	if tx != nil {
		return fc()
	}

	db := client_singleton.DBs[label]
	if db == nil {
		panic(fmt.Sprintf("Can't get DB, DB label:%s", label))
	}

	isPanic := true
	tx = db.Begin()
	c.txs[label] = tx
	defer func() {
		delete(c.txs, label)

		if isPanic || err != nil {
			if dbErr := tx.Rollback().Error; dbErr != nil {
				err = err_code.WriteDBError.WithErr(dbErr)
			}
			return
		}

		if dbErr := tx.Commit().Error; dbErr != nil {
			err = err_code.WriteDBError.WithErr(dbErr)
		}
	}()

	err = fc()
	isPanic = false
	return err
}
