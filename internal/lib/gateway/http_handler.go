package gateway

import (
	"github.com/Me1onRind/go-demo/internal/lib/ctm_context"
	"github.com/Me1onRind/go-demo/internal/lib/err_code"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type HTTPHandler func(c *ctm_context.Context, raw interface{}) (data interface{}, err *err_code.Error)

func JSON(handler HTTPHandler, paramType interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		var err *err_code.Error
		var data interface{}
		var response *JSONResponse
		var requestParams interface{}

		commonCtx := ctm_context.GetCtmCtxFromGinCtx(c)

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
