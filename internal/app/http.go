package app

import (
	"github.com/Me1onRind/go-demo/internal/infrastructure/middleware"
	"github.com/Me1onRind/go-demo/internal/usecase/unexpectuc"
	"github.com/Me1onRind/go-demo/internal/usecase/useruc"
	"github.com/gin-gonic/gin"
)

type HttpServer struct {
	UserUsecase     *useruc.UserUsecase
	UnexpectUsecase *unexpectuc.UnexcpectUsecase
}

func NewHttpServer() *HttpServer {
	h := &HttpServer{
		UserUsecase:     useruc.NewUserUsecase(),
		UnexpectUsecase: unexpectuc.NewUnexpectUseCase(),
	}
	return h
}

func (h *HttpServer) RegisterMiddleware(r *gin.Engine) *HttpServer {
	r.Use(
		middleware.Recover(),
		middleware.SetRequestId(),
		middleware.AccessLog(),
	)
	return h
}

func (h *HttpServer) ReigsterRouter(router *gin.RouterGroup) {
	router = router.Group("/api")
	userGroup := router.Group("/user")
	userGroup.GET("get_user_detail", middleware.JSON(h.UserUsecase.QueryUserInfo, nil))

	unexpectGroup := router.Group("/unexpect")
	unexpectGroup.GET("/panic", middleware.JSON(h.UnexpectUsecase.Panic, nil))
}
