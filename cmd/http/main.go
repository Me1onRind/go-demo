package main

import (
	"github.com/Me1onRind/go-demo/internal/lib/initialize"
	"github.com/Me1onRind/go-demo/internal/lib/middleware"
	"github.com/Me1onRind/go-demo/internal/router"
	"github.com/gin-gonic/gin"
)

func Init() {
	funcs := []func() error{
		initialize.InitLogger,
		initialize.InitOpentracking("go-http-demo", "0.0.1"),
		initialize.InitGrpcClients,
	}

	for _, v := range funcs {
		if err := v(); err != nil {
			panic(err)
		}
	}
}

func main() {
	Init()
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

	if err := r.Run("0.0.0.0:8081"); err != nil {
		panic(err)
	}
}
