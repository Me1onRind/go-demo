package middleware

import (
	"bytes"
	"context"
	"io/ioutil"
	"time"

	"github.com/Me1onRind/go-demo/constant/sys_constant"
	"github.com/Me1onRind/go-demo/global/logger_singleton"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/Me1onRind/go-demo/infrastructure/logger"
)

func GrpcAccessLog() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		begin := time.Now()
		defer func() {
			logger.CtxInfo(ctx, "access request", zap.Reflect("req", req), zap.Reflect("resp", resp),
				zap.String("method", info.FullMethod), zap.Error(err), zap.Duration("cost", time.Since(begin)),
			)
		}()
		resp, err = handler(ctx, req)
		return resp, err
	}
}

func GinSetContextLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceId := c.GetString(sys_constant.TraceIdKey)
		loggerInstance := logger_singleton.Logger.With(
			zap.String("traceId", traceId),
		)
		c.Set(sys_constant.LoggerKey, loggerInstance)
	}
}

func GrpcSetContextLogger() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		value := ctx.Value(sys_constant.TraceIdKey)
		if value != nil {
			if traceId, ok := value.(string); ok {
				loggerInstance := logger_singleton.Logger.With(
					zap.String("traceId", traceId),
				)
				ctx = context.WithValue(ctx, sys_constant.LoggerKey, loggerInstance)
			}
		}
		return handler(ctx, req)
	}
}

func GinAccessLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request []byte

		contentType := c.ContentType()
		if contentType == "application/json" || contentType == "text/plain" {
			var err error
			request, err = c.GetRawData()
			if err != nil {
				logger.CtxError(c, "Get request body error", zap.Error(err))
			}
			c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(request))
		}

		lw := &logWriter{
			ResponseWriter: c.Writer,
			buff:           &bytes.Buffer{},
		}
		c.Writer = lw

		start := time.Now()
		defer func() {
			end := time.Now()
			logger.CtxInfo(c, "access request",
				zap.String("proto", c.Request.Proto), zap.String("method", c.Request.Method), zap.String("path", c.Request.URL.Path), zap.String("rawQuery", c.Request.URL.RawQuery),
				zap.String("reqBody", string(request)), zap.String("clientIP", c.ClientIP()), zap.String("resp", lw.buff.String()), zap.Duration("cost", end.Sub(start)),
			)
		}()

		c.Next()

	}
}

type logWriter struct {
	gin.ResponseWriter
	buff *bytes.Buffer
}

func (w *logWriter) Write(b []byte) (int, error) {
	w.buff.Write(b)
	return w.ResponseWriter.Write(b)
}
