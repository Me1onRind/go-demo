package unit_test

import (
	"os"

	"github.com/Me1onRind/go-demo/config"
	"github.com/Me1onRind/go-demo/global/client_singleton"
	"github.com/Me1onRind/go-demo/global/logger_singleton"
	"github.com/Me1onRind/go-demo/infrastructure/client/redis_client"
	"github.com/alicebob/miniredis/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitGlobalVar() {
	r, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	client_singleton.RedisClient = redis_client.NewRedisClient(&config.RedisConfig{
		Addr: r.Addr(),
	})

	encoderCfg := zapcore.EncoderConfig{}
	core := zapcore.NewCore(zapcore.NewJSONEncoder(encoderCfg), os.Stdout, zap.PanicLevel)
	logger_singleton.Logger = zap.New(core)
}
