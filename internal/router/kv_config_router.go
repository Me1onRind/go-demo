package router

//import (
//"github.com/Me1onRind/go-demo/internal/controller/kv_config_controller"
//"github.com/Me1onRind/go-demo/internal/lib/gateway"
//protocol "github.com/Me1onRind/go-demo/protocol/http_proto/kv_config_protocol"
//"github.com/gin-gonic/gin"
//)

//func SetKVConfigRouter(router gin.IRouter) {
//group := router.Group("/kv_config")
//kvConfigController := kv_config_controller.NewKvConfigController()
//group.POST("create", gateway.JSON(kvConfigController.CreateKvConfig, &protocol.CreateKVconfigReq{}))
//group.POST("update", gateway.JSON(kvConfigController.UpdateKvConfig, &protocol.UpdateKVconfigReq{}))
//group.GET("get", gateway.JSON(kvConfigController.GetKvConfigByID, &protocol.GetKVconfigByIDReq{}))
//}
