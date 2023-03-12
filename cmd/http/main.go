package main

import (
	"context"
	"fmt"

	"github.com/Me1onRind/go-demo/internal/app"
	"github.com/Me1onRind/go-demo/internal/infrastructure/goroutine"
	"github.com/Me1onRind/go-demo/internal/infrastructure/logger"
	"github.com/gin-gonic/gin"
)

func main() {
	goroutine.SafeGo(context.Background(), func() {
		_ = app.NewMetrics().Run()
	})
	r := gin.New()

	app.NewHttpServer().
		RegisterMiddleware(r).
		RegisterRouter(r.Group("/")).
		Init()

	logger.Infof("Server run")
	if err := r.Run(); err != nil {
		fmt.Println(err)
	}
}
