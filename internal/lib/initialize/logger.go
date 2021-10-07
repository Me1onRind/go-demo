package initialize

import (
	"os"

	"github.com/Me1onRind/go-demo/internal/core/logger"
	"github.com/opentracing/opentracing-go"
	"go.elastic.co/apm"
	"go.elastic.co/apm/module/apmot"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitLogger() error {
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
	logger.Logger = zap.New(core, zap.AddCaller())

	return nil
}

func CloseLogger() error {
	return nil
	//return os.Stdout.Close()
}

func InitOpentracking(serviceName, version string) func() error {
	return func() error {
		os.Setenv("ELASTIC_APM_SERVER_URL", "http://localhost:8200")
		os.Setenv("ELASTIC_APM_SECRET_TOKEN", "")
		os.Setenv("ELASTIC_APM_STACK_TRACE_LIMIT", "0")
		os.Setenv("ELASTIC_APM_USE_ELASTIC_TRACEPARENT_HEADER", "false")
		os.Setenv("ELASTIC_APM_TRANSACTION_SAMPLE_RATE", "1.0")
		tracer, err := apm.NewTracer(serviceName, version)
		if err != nil {
			return err
		}
		opentracing.SetGlobalTracer(apmot.New(apmot.WithTracer(tracer)))
		return nil
	}
}
