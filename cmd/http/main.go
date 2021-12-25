package main

import (
	"context"

	"github.com/Me1onRind/go-demo/infrastructure/initialize"
	"github.com/Me1onRind/go-demo/infrastructure/middleware"
	"github.com/Me1onRind/go-demo/internal/router"
	"github.com/gin-gonic/gin"
)

func Init(ctx context.Context) {
	funcs := []func() error{
		initialize.InitLogger(),
		initialize.InitOpentracking("go-http-demo", "0.0.1"),
		//initialize.InitLocalConfig("./conf"),
		//initialize.InitEtcdClient,
		//initialize.InitEtcdConfig(ctx, "/go-demo/config.yml"),
		//initialize.InitMysqlClients,
		//initialize.InitRedisClient,
		////initialize.InitLocalCache,
		//initialize.InitGrpcClients,
	}

	for _, v := range funcs {
		if err := v(); err != nil {
			panic(err)
		}
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	Init(ctx)

	defer func() {
		cancel()
	}()

	r := gin.Default()
	apiGroup := r.Group("/api")
	apiGroup.Use(
		middleware.GinRecover(),
		middleware.GinTracer(),
		middleware.GinSetContextLogger(),
		middleware.GinAccessLog(),
		//middleware.GinContext(),
		//middleware.GinRecover(),
	//middleware.GinLogger(),
	)

	router.SetFooRouter(apiGroup)
	//router.SetPeridoicTaskRouter(apiGroup)
	//router.SetKVConfigRouter(apiGroup)

	if err := r.Run("0.0.0.0:8081"); err != nil {
		panic(err)
	}
}
