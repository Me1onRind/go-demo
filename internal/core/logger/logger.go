package logger

import (
	"go.uber.org/zap"
)

var (
	Logger       *zap.Logger
	StdoutLogger *zap.Logger
)

//type AsynqLogger struct {
//}

//func NewAsyncLogger() *AsynqLogger {
//a := &AsynqLogger{}
//return a
//}

//func (a *AsynqLogger) Debug(args ...interface{}) {
//}
