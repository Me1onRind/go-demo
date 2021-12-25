package gateway

import (
	"github.com/Me1onRind/go-demo/constant/sys_constant"
	"github.com/Me1onRind/go-demo/err_code"
	"github.com/Me1onRind/go-demo/infrastructure/ctm_context"
	"github.com/Me1onRind/go-demo/infrastructure/logger"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

type HTTPHandler func(c *ctm_context.Context, raw interface{}) (data interface{}, err *err_code.Error)

func JSON(handler HTTPHandler, paramType interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		var err *err_code.Error
		var data interface{}
		var response *JSONResponse
		var requestParams interface{}

		defer func() {
			if response != nil {
				value, exist := c.Get(sys_constant.SpanKey)
				if exist {
					if span, ok := value.(opentracing.Span); ok {
						span.SetTag("errcode", response.Errcode)
						span.SetTag("message", response.Message)
					}
				}
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

		commonCtx := ctm_context.NewContext(c)
		data, err = handler(commonCtx, requestParams)
		logger.CtxInfo(commonCtx, "json gateway", zap.Reflect("inputObject", requestParams), zap.Reflect("outputObject", data), zap.Reflect("error", err))
		if err == nil {
			err = err_code.SUCCESS
		}
		response = NewJSONResponse(err, data)
	}
}
