package main

import (
	"context"

	"github.com/Me1onRind/go-demo/internal/core/logger"
	"github.com/Me1onRind/go-demo/internal/lib/initialize"
	"github.com/Me1onRind/go-demo/internal/lib/middleware"
	"github.com/Me1onRind/go-demo/internal/router"
	"github.com/gin-gonic/gin"
)

func Init(ctx context.Context) {
	funcs := []func() error{
		initialize.InitLogger,
		initialize.InitLocalConfig("./conf"),
		initialize.InitEtcdClient,
		initialize.InitEtcdConfig(ctx, "/go-demo/config.yml"),
		initialize.InitOpentracking("go-http-demo", "0.0.1"),
		initialize.InitMysqlClients,
		initialize.InitRedisClient,
		//initialize.InitLocalCache,
		initialize.InitGrpcClients,
	}

	for _, v := range funcs {
		if err := v(); err != nil {
			panic(err)
		}
	}
}

func Close() {
	logger.StdoutLogger.Info("Process exit close")

	funcs := []func() error{
		initialize.CloseLogger,
		initialize.CloseEtcdClient,
		initialize.CloseMysqlClients,
		initialize.CloseRedisClient,
		initialize.CloseGrpcClients,
	}

	for _, v := range funcs {
		if err := v(); err != nil {
			logger.StdoutLogger.Error("Process exit close")
		}
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	Init(ctx)

	defer func() {
		cancel()
		Close()
	}()

	r := gin.Default()
	apiGroup := r.Group("/api")
	apiGroup.Use(
		middleware.GinContext(),
		middleware.GinRecover(),
		middleware.GinTracer(),
		middleware.GinLogger(),
	)

	router.SetFooRouter(apiGroup)
	router.SetPeridoicTaskRouter(apiGroup)
	router.SetKVConfigRouter(apiGroup)

	if err := r.Run("0.0.0.0:8081"); err != nil {
		panic(err)
	}
}
