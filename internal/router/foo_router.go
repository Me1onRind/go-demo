package router

import (
	"github.com/Me1onRind/go-demo/internal/controller/foo_controller"
	"github.com/Me1onRind/go-demo/internal/lib/gateway"
	"github.com/Me1onRind/go-demo/protocol"
	"github.com/gin-gonic/gin"
)

func SetFooRouter(router gin.IRouter) {
	group := router.Group("/foo")
	fooController := foo_controller.NewFooController()
	group.POST("greet", gateway.JSON(fooController.ProxyGreet, &protocol.GreetProxyRequest{}))
}
