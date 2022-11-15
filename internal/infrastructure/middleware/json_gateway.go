package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Me1onRind/go-demo/constant/code"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
)

type JsonResponse struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type HTTPHandler func(c context.Context, raw any) (data any, err error)

func JSON(handler HTTPHandler, paramType any) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := mustGetGinExtractContext(c)
		data, err := handler(ctx, nil)
		if err != nil {
		}

		jsonData, err := jsoniter.Marshal(&JsonResponse{
			Code:    code.Success,
			Message: "success",
			Data:    data,
		})
		if err != nil {
			jsonData, _ = jsoniter.Marshal(&JsonResponse{
				Code:    code.JsonEncodeFailed,
				Message: fmt.Sprintf("JSON Gateway encode response failed, err:[%s]", err.Error()),
			})
		}

		c.Data(http.StatusOK, "application/json; charset=utf-8", jsonData)
		c.Next()
	}
}
