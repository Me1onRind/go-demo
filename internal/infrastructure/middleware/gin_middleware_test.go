package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Me1onRind/go-demo/internal/infrastructure/tracer"
	"github.com/gin-gonic/gin"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/stretchr/testify/assert"
	"go.elastic.co/apm"
	"go.elastic.co/apm/module/apmhttp"
	"go.elastic.co/apm/module/apmot"
)

func Test_Recover(t *testing.T) {
	c, r := gin.CreateTestContext(httptest.NewRecorder())
	r.Use(Recover())
	r.POST("/panic", func(ctx *gin.Context) {
		panic("test panic")
	})
	c.Request, _ = http.NewRequest("POST", "/panic", nil)
	assert.NotPanics(t, func() { r.HandleContext(c) })
}

func Test_Tracer(t *testing.T) {
	t.Run("Root_Tracer", func(t *testing.T) {
		tr, _ := apm.NewTracer("test", "0.0.1")
		opentracing.SetGlobalTracer(apmot.New(apmot.WithTracer(tr)))
		c, r := gin.CreateTestContext(httptest.NewRecorder())
		r.Use(Tracer())
		r.POST("/tracer", func(ctx *gin.Context) {
		})
		c.Request, _ = http.NewRequest("POST", "/tracer", nil)
		r.HandleContext(c)
		assert.NotEmpty(t, tracer.GetTraceId(mustGetGinExtractContext(c)))
		assert.NotEmpty(t, tracer.GetSpanId(mustGetGinExtractContext(c)))
	})

	t.Run("Children_Tracer", func(t *testing.T) {
		tr, _ := apm.NewTracer("test", "0.0.1")
		opentracing.SetGlobalTracer(apmot.New(apmot.WithTracer(tr)))
		c, r := gin.CreateTestContext(httptest.NewRecorder())
		r.Use(Tracer())
		r.POST("/tracer", func(ctx *gin.Context) {
		})
		c.Request, _ = http.NewRequest("POST", "/tracer", nil)
		c.Request.Header[apmhttp.W3CTraceparentHeader] = []string{"00-6d4cea0cef8ad69bc2adde7ed03afcba-6d4cea0cef8ad69b-01"}
		r.HandleContext(c)
		assert.Equal(t, "6d4cea0cef8ad69bc2adde7ed03afcba", tracer.GetTraceId(mustGetGinExtractContext(c)))
		assert.Equal(t, "6d4cea0cef8ad69b", tracer.GetSpanId(mustGetGinExtractContext(c)))
	})

	t.Run("Root_Tracer_Wrong_Http_Header", func(t *testing.T) {
		tr, _ := apm.NewTracer("test", "0.0.1")
		opentracing.SetGlobalTracer(apmot.New(apmot.WithTracer(tr)))
		c, r := gin.CreateTestContext(httptest.NewRecorder())
		r.Use(Tracer())
		r.POST("/tracer", func(ctx *gin.Context) {
		})
		c.Request, _ = http.NewRequest("POST", "/tracer", nil)
		c.Request.Header[apmhttp.W3CTraceparentHeader] = []string{"6d4cea0cef8ad69bc2adde7ed03afcba-6d4cea0cef8ad69b-01"}
		r.HandleContext(c)
		assert.NotEmpty(t, tracer.GetTraceId(mustGetGinExtractContext(c)))
		assert.NotEmpty(t, tracer.GetSpanId(mustGetGinExtractContext(c)))
	})

	t.Run("Root_Tracer_Empty_Http_Header", func(t *testing.T) {
		tr, _ := apm.NewTracer("test", "0.0.1")
		opentracing.SetGlobalTracer(apmot.New(apmot.WithTracer(tr)))
		c, r := gin.CreateTestContext(httptest.NewRecorder())
		r.Use(Tracer())
		r.POST("/tracer", func(ctx *gin.Context) {
		})
		c.Request, _ = http.NewRequest("POST", "/tracer", nil)
		r.HandleContext(c)
		assert.NotEmpty(t, tracer.GetTraceId(mustGetGinExtractContext(c)))
		assert.NotEmpty(t, tracer.GetSpanId(mustGetGinExtractContext(c)))
	})
}

func Test_AccessLog(t *testing.T) {
	c, r := gin.CreateTestContext(httptest.NewRecorder())
	r.Use(AccessLog())
	r.POST("/logger", func(ctx *gin.Context) {
		ctx.JSON(200, map[string]any{
			"code": 100,
		})
	})
	c.Request, _ = http.NewRequest("POST", "/logger", strings.NewReader(`{"key":"value"}`))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.Host = "127.0.0.1"
	assert.NotPanics(t, func() {
		r.HandleContext(c)
	})
}
