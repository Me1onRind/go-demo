package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Me1onRind/go-demo/internal/infrastructure/logger"
	customErr "github.com/Me1onRind/go-demo/internal/model/errors"
	"github.com/Me1onRind/go-demo/protocol/code"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
)

type JsonResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type HTTPHandler[T any] func(c context.Context, request *T) (data any, err error)

func JSON[T any](handler HTTPHandler[T]) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Data(http.StatusOK, "application/json; charset=utf-8", jsonGateWay(c, handler))
		c.Next()
	}
}

func jsonGateWay[T any](c *gin.Context, handler HTTPHandler[T]) []byte {
	ctx := mustGetGinExtractContext(c)

	var response *JsonResponse
	var request T
	if err := c.ShouldBind(&request); err != nil {
		response = &JsonResponse{
			Code:    code.ProtocolDecodeFail,
			Message: fmt.Sprintf("Decode request text fail, cause:[%s]", err),
		}
	} else {
		data, err := handler(ctx, &request)
		response = getResponse(data, err)
	}

	jsonData, err := jsoniter.Marshal(response)
	if err != nil {
		logger.CtxErrorf(ctx, "Marshal response fail, err:[%s]", err)
		jsonData, _ = jsoniter.Marshal(&JsonResponse{
			Code:    code.JsonEncodeFail,
			Message: fmt.Sprintf("JSON Gateway encode response fail, err:[%s]", err.Error()),
		})
	}
	return jsonData
}

func getResponse(data any, err error) *JsonResponse {
	response := &JsonResponse{}
	if err == nil {
		response.Code = code.Success
		response.Message = "Success"
		response.Data = data
		return response
	}
	response.Message = err.Error()
	if expectErr := customErr.ExtractError(err); expectErr != nil {
		response.Code = expectErr.Code
	} else {
		response.Code = code.Unexpect
	}

	if code.IsWarning(response.Code) {
		response.Data = data
	}
	return response
}
