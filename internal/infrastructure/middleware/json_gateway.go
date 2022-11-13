package middleware

import (
	"net/http"

	"github.com/Me1onRind/go-demo/internal/infrastructure/base"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
)

type HTTPHandler func(c *base.Context, raw interface{}) (data interface{}, err error)

func JSON(handler HTTPHandler, paramType interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		data, err := handler(nil, nil)
		if err != nil {
		}

		jsonData, err := jsoniter.Marshal(data)
		if err != nil {
		}

		c.Data(http.StatusOK, "application/json; charset=utf-8", jsonData)
		c.Next()
	}
}
