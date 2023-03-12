package logger

import (
	"context"
	"time"

	"gorm.io/gorm/logger"
)

func NewMysqlLogger() logger.Interface {
	return &mysqlLogger{}
}

type mysqlLogger struct {
}

func (m *mysqlLogger) LogMode(logLevel logger.LogLevel) logger.Interface {
	return m
}

func (m *mysqlLogger) Info(ctx context.Context, format string, a ...interface{}) {
	globalLogger.Infof(format, a...)
}

func (m *mysqlLogger) Warn(ctx context.Context, format string, a ...interface{}) {
	getLoggerFromCtx(ctx).Warnf(format, a...)
}

func (m *mysqlLogger) Error(ctx context.Context, format string, a ...interface{}) {
	getLoggerFromCtx(ctx).Errorf(format, a...)
}

func (m *mysqlLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	sql, rowsAffected := fc()
	if err != nil {
		getLoggerFromCtx(ctx).Errorf("rows:[%d],sql:[%s],duration:[%s],error:[%s]", rowsAffected, sql, time.Since(begin), err)
	} else {
		getLoggerFromCtx(ctx).Infof("rows:[%d],sql:[%s],duration:[%s]", rowsAffected, sql, time.Since(begin))
	}
}
