package gateway

import (
	"github.com/Me1onRind/go-demo/internal/core/common"
	"github.com/Me1onRind/go-demo/internal/err_code"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type HTTPHandler func(c *common.Context, raw interface{}) (data interface{}, err *common.Error)

func JSON(handler HTTPHandler, paramType interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		var err *common.Error
		var data interface{}
		var response *JSONResponse
		var requestParams interface{}

		commonCtx := common.GetContext(c)

		defer func() {
			if response != nil {
				span := commonCtx.Span
				span.SetTag("errcode", response.Errcode)
				span.SetTag("message", response.Message)
				c.JSON(200, response)
			}
		}()

		if paramType != nil {
			requestParams = parserProtocol(paramType)
			if e := c.ShouldBind(requestParams); e != nil {
				response = NewJSONResponse(err_code.InvalidParamError.WithErr(e), nil)
				return
			}
		}

		data, err = handler(commonCtx, requestParams)
		commonCtx.Logger.Info("json gateway", zap.Reflect("inputObject", requestParams), zap.Reflect("outputObject", data), zap.Reflect("error", err))
		if err == nil {
			err = err_code.SUCCESS
		}
		response = NewJSONResponse(err, data)
	}
}
