package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"reflect"

	"github.com/Me1onRind/go-demo/internal/infrastructure/logger"
	customErr "github.com/Me1onRind/go-demo/internal/model/errors"
	"github.com/Me1onRind/go-demo/protocol/code"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
)

type JsonResponse struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type HTTPHandler func(c context.Context, raw any) (data any, err error)

func JSON(handler HTTPHandler, protocol any) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Data(http.StatusOK, "application/json; charset=utf-8", jsonGateWay(c, handler, protocol))
		c.Next()
	}
}

func jsonGateWay(c *gin.Context, handler HTTPHandler, protocol any) []byte {
	ctx := mustGetGinExtractContext(c)

	var response *JsonResponse
	raw, err := initProtocol(c, protocol)

	if err != nil {
		response = &JsonResponse{
			Code:    code.ProtocolDecodeFail,
			Message: fmt.Sprintf("Decode request text fail, cause:[%s]", err),
		}
	} else {
		data, err := handler(ctx, raw)
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
	var expectErr *customErr.Error
	if errors.As(err, &expectErr) {
		response.Code = expectErr.Code
	} else {
		response.Code = code.Unexpect
	}

	if code.IsWarning(response.Code) {
		response.Data = data
	}
	return response
}

func initProtocol(c *gin.Context, protocol any) (any, error) {
	if protocol == nil {
		return nil, nil
	}
	entity, err := newProtocol(protocol)
	if err != nil {
		return nil, err
	}
	if err := c.ShouldBind(entity); err != nil {
		return nil, err
	}

	return entity, nil
}

func newProtocol(protocol any) (any, error) {
	valueType := reflect.TypeOf(protocol)
	switch valueType.Kind() {
	case reflect.Ptr:
		return reflect.New(valueType.Elem()).Interface(), nil
	case reflect.Struct:
		return reflect.New(valueType).Interface(), nil
	default:
		return nil, fmt.Errorf("New protocol:[%+v] struct fail, cause it's type:[%s] not support", protocol, valueType.Kind())
	}
}
