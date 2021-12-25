package router

import (
	"github.com/Me1onRind/go-demo/infrastructure/gateway"
	"github.com/Me1onRind/go-demo/internal/controller/foo_controller"
	"github.com/Me1onRind/go-demo/protocol/http_proto"
	"github.com/gin-gonic/gin"
)

func SetFooRouter(router gin.IRouter) {
	group := router.Group("/foo")
	fooController := foo_controller.NewFooController()
	group.POST("greet", gateway.JSON(fooController.ProxyGreet, &http_proto.GreetProxyRequest{}))
}
