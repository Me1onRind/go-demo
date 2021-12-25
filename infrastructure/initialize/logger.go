package initialize

import (
	"os"

	"github.com/Me1onRind/go-demo/global/logger_singleton"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitLogger() func() error {
	return func() error {
		config := zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.EpochTimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		}

		config.EncodeDuration = zapcore.MillisDurationEncoder
		config.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
		config.EncodeLevel = zapcore.CapitalLevelEncoder

		core := zapcore.NewCore(zapcore.NewJSONEncoder(config), zapcore.AddSync(os.Stdout), zapcore.InfoLevel)
		logger_singleton.Logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

		return nil
	}
}
