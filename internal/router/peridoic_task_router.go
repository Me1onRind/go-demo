package router

//import (
//"github.com/Me1onRind/go-demo/internal/controller/peridoic_task_controller"
//"github.com/Me1onRind/go-demo/internal/lib/gateway"
//protocol "github.com/Me1onRind/go-demo/protocol/http_proto/peridoic_task_protocol"
//"github.com/gin-gonic/gin"
//)

//func SetPeridoicTaskRouter(router gin.IRouter) {
//group := router.Group("/peridoic_task")
//peridoicTaskController := peridoic_task_controller.NewPeridoicTaskController()
//group.POST("create", gateway.JSON(peridoicTaskController.CreatePeridoicTask, &protocol.CreatePeridoicTaskReq{}))
//group.POST("update", gateway.JSON(peridoicTaskController.UpdatePeridoicTask, &protocol.UpdatePeridoicTaskReq{}))
//group.GET("get", gateway.JSON(peridoicTaskController.GetPeridoicTaskByID, &protocol.GetPeridoicTaskByIDReq{}))
//}
