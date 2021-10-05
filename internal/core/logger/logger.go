package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	Logger       *zap.Logger
	StdoutLogger *zap.Logger
)

func init() {
	stdoutConfig := zap.NewProductionEncoderConfig()
	stdoutConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	stdoutConfig.EncodeDuration = zapcore.MillisDurationEncoder
	stdoutConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
	stdoutConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	stdoutConfig.ConsoleSeparator = "|"
	stdoutCore := zapcore.NewCore(zapcore.NewConsoleEncoder(stdoutConfig), zapcore.AddSync(os.Stdout), zapcore.InfoLevel)
	StdoutLogger = zap.New(stdoutCore, zap.AddCaller())
}
