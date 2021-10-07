package ctm_context

import (
	"context"
	"fmt"

	"github.com/Me1onRind/go-demo/global/store"
	"github.com/Me1onRind/go-demo/global/store/db_label"
	"github.com/Me1onRind/go-demo/internal/core/logger"
	"github.com/Me1onRind/go-demo/internal/lib/err_code"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const (
	ctmCtxKey = "ctm"
)

type Context struct {
	context.Context

	Logger *zap.Logger
	Span   opentracing.Span

	txs map[db_label.Label]*gorm.DB
}

func NewContext(ctx context.Context) *Context {
	c := &Context{
		Context: ctx,
		Logger:  logger.Logger,
		txs:     map[db_label.Label]*gorm.DB{},
	}
	return c
}

func (c *Context) DB(label db_label.Label) *gorm.DB {
	tx := c.txs[label]
	if tx != nil {
		return tx
	}

	db := store.DBs[label]
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

	db := store.DBs[label]
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

func (ctx *Context) SaveInGinCtx(c *gin.Context) {
	c.Set(ctmCtxKey, ctx)
}

func GetCtmCtxFromGinCtx(c *gin.Context) *Context {
	ctx, _ := c.Get(ctmCtxKey)
	return ctx.(*Context)
}

func ContextLogger(ctx context.Context) *zap.Logger {
	if commonCtx, ok := ctx.(*Context); ok {
		return commonCtx.Logger
	}

	return logger.StdoutLogger
}
