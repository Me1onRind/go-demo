package app

import (
	"context"

	"github.com/Me1onRind/go-demo/internal/infrastructure/logger"
	"github.com/Me1onRind/go-demo/internal/infrastructure/middleware"
	"github.com/Me1onRind/go-demo/internal/initialize"
	"github.com/Me1onRind/go-demo/internal/usecase/unexpectuc"
	"github.com/Me1onRind/go-demo/internal/usecase/useruc"
	"github.com/Me1onRind/go-demo/protocol/userproto"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
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

func (h *HttpServer) ReigsterRouter(router *gin.RouterGroup) *HttpServer {
	router = router.Group("/api")
	userGroup := router.Group("/user")
	userGroup.GET("get_user_detail", middleware.JSON(h.UserUsecase.GetUserDetail, &userproto.GetUserDetailReq{}))
	userGroup.POST("create_user", middleware.JSON(h.UserUsecase.CreateUser, &userproto.CreateUserReq{}))

	unexpectGroup := router.Group("/unexpect")
	unexpectGroup.GET("/panic", middleware.JSON(h.UnexpectUsecase.Panic, nil))
	return h
}

func (h *HttpServer) Init() *HttpServer {
	initFuncs := []initialize.InitHandler{
		initialize.InitFileConfig("./conf.yml"),
		initialize.InitEtcdClient(),
		initialize.InitDynamicConfig(),
		initialize.InitMysqlClient(),
	}
	ctx := context.Background()
	ctx = logger.WithFields(ctx, logrus.Fields{
		logger.RequestIdKey: "main-goruntine",
	})

	for _, f := range initFuncs {
		if err := f(ctx); err != nil {
			logger.CtxFatalf(ctx, "initialize fail, err:[%s]", err)
		}
	}
	return h
}
