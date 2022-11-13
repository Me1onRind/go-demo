package app

import (
	"github.com/Me1onRind/go-demo/internal/infrastructure/middleware"
	"github.com/Me1onRind/go-demo/internal/usecase/userus"
	"github.com/gin-gonic/gin"
)

type HttpServer struct {
	UserUsecase *userus.UserUsecase
}

func NewHttpServer() *HttpServer {
	h := &HttpServer{
		UserUsecase: userus.NewUserUsecase(),
	}
	return h
}

func (h *HttpServer) RegisterMiddleware(r *gin.Engine) *HttpServer {
	r.Use(
		middleware.ExtractRequestCtx(),
		middleware.SetRequestId(),
		middleware.AccessLog(),
	)
	return h
}

func (h *HttpServer) ReigsterRouter(router *gin.RouterGroup) {
	router = router.Group("/api")
	userGroup := router.Group("/user")
	userGroup.GET("get_user_detail", middleware.JSON(h.UserUsecase.QueryUserInfo, nil))
}
