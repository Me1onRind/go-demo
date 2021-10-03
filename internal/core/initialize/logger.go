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
	config := zap.NewProductionEncoderConfig()
	config.EncodeDuration = zapcore.MillisDurationEncoder
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	//f, err := os.OpenFile("./log/info.log", os.O_WRONLY|os.O_RDONLY|os.O_APPEND|os.O_CREATE, 0755)
	//if err != nil {
	//panic(err)
	//}
	core := zapcore.NewCore(zapcore.NewJSONEncoder(config), zapcore.AddSync(os.Stdout), zapcore.InfoLevel)
	logger.Logger = zap.New(core, zap.AddCaller())
	return nil
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
