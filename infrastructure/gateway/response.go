package gateway

import (
	"github.com/Me1onRind/go-demo/err_code"
)

type JSONResponse struct {
	Errcode int32       `json:"errcode"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func NewJSONResponse(err *err_code.Error, data interface{}) *JSONResponse {
	if data == nil {
		data = struct{}{}
	}
	return &JSONResponse{
		Errcode: err.Code,
		Message: err.String(),
		Data:    data,
	}
}
