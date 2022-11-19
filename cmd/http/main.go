package main

import (
	"fmt"

	"github.com/Me1onRind/go-demo/internal/app"
	"github.com/Me1onRind/go-demo/internal/infrastructure/logger"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.New()

	app.NewHttpServer().
		RegisterMiddleware(r).
		ReigsterRouter(r.Group("/")).
		Init()

	logger.Infof("Server run")
	if err := r.Run(); err != nil {
		fmt.Println(err)
	}
}
