package middleware

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/Me1onRind/go-demo/internal/global/gerror"
	"github.com/Me1onRind/go-demo/protocol"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func Test_jsonGateWay_Empty(t *testing.T) {
	tests := []struct {
		name   string
		method string
		body   string
		param  string
		output string
	}{
		{
			name:   "POST",
			method: "POST",
			body:   `{"a":1123}`,
			output: `{"code":0,"message":"Success","data":null}`,
		},
		{
			name:   "GET",
			body:   `a=123`,
			output: `{"code":0,"message":"Success","data":null}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request, _ := http.NewRequest(test.method, "http://localhost/index?"+test.param, strings.NewReader(test.body))
			ctx := &gin.Context{Request: request}
			res := jsonGateWay(ctx, func(c context.Context, r *protocol.EmptyReq) (any, error) {
				return nil, nil
			})
			assert.Equal(t, string(res), test.output)
		})
	}
}

func Test_jsonGateWay_Value(t *testing.T) {
	tests := []struct {
		name   string
		method string
		body   string
		param  string
		output string
	}{
		{
			name:   "POST",
			method: "POST",
			body:   `{"a":123}`,
			output: `{"code":0,"message":"Success","data":{"a":123}}`,
		},
		{
			name:   "GET",
			param:  `a=123`,
			output: `{"code":0,"message":"Success","data":{"a":123}}`,
		},
		{
			name:   "POST_Fail",
			method: "POST",
			body:   `{"b":123}`,
			output: `{"code":-100003,"message":"Decode request text fail, cause:[Key: 'Req.A' Error:Field validation for 'A' failed on the 'required' tag]","data":null}`,
		},
		{
			name:   "GET_Fail",
			method: "GET",
			body:   `{"b":123}`,
			output: `{"code":-100003,"message":"Decode request text fail, cause:[Key: 'Req.A' Error:Field validation for 'A' failed on the 'required' tag]","data":null}`,
		},
	}

	type Req struct {
		A int `form:"a" json:"a" binding:"required"`
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request, _ := http.NewRequest(test.method, "http://localhost/index?"+test.param, strings.NewReader(test.body))
			if test.method == "POST" {
				request.Header["Content-Type"] = []string{"application/json"}
			}
			ctx := &gin.Context{Request: request}
			res := jsonGateWay(ctx, func(c context.Context, r *Req) (any, error) {
				return map[string]any{
					"a": r.A,
				}, nil
			})
			assert.Equal(t, string(res), test.output)
		})
	}
}

func Test_jsonGateWay_Marshal_Fail(t *testing.T) {
	request, _ := http.NewRequest("GET", "http://localhost/index", nil)
	ctx := &gin.Context{Request: request}
	res := jsonGateWay(ctx, func(c context.Context, r *protocol.EmptyReq) (any, error) {
		return func() {}, nil
	})
	assert.Equal(t, string(res), `{"code":-100001,"message":"JSON Gateway encode response fail, err:[middleware.JsonResponse.Data: func() is unsupported type]","data":null}`)
}

func Test_jsonGateWay_Normal_Fail(t *testing.T) {
	tests := []struct {
		name   string
		data   any
		err    error
		output string
	}{
		{
			name:   "emptyData",
			err:    io.EOF,
			output: `{"code":-100000,"message":"EOF","data":null}`,
		},
		{
			name:   "hasData",
			err:    io.EOF,
			data:   map[string]any{"key": "value"},
			output: `{"code":-100000,"message":"EOF","data":null}`,
		},
		{
			name:   "customError",
			err:    gerror.ReadDBError,
			output: `{"code":-100004,"message":"Read DB Fail","data":null}`,
		},
		{
			name:   "warning",
			data:   map[string]any{"key": "value"},
			err:    gerror.DuplicateError,
			output: `{"code":1,"message":"Duplicate Request","data":{"key":"value"}}`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request, _ := http.NewRequest("GET", "http://localhost/index", nil)
			ctx := &gin.Context{Request: request}
			res := jsonGateWay(ctx, func(c context.Context, r *protocol.EmptyReq) (any, error) {
				return test.data, test.err
			})
			assert.Equal(t, string(res), test.output)
		})
	}
}
