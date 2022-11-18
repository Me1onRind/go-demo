package middleware

import (
	"bytes"
	"io"
	"time"

	"github.com/Me1onRind/go-demo/internal/infrastructure/logger"
	"github.com/gin-gonic/gin"
)

func AccessLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			request []byte
			err     error
		)

		contentType := c.ContentType()
		if contentType == "application/json" || contentType == "text/plain" {
			request, err = c.GetRawData()
			if err != nil {
				logger.CtxErrorf(mustGetGinExtractContext(c), "GetRawData faile, cause:[%s]", err)
			} else {
				c.Request.Body = io.NopCloser(bytes.NewBuffer(request))
			}
		}
		lw := &logWriter{
			ResponseWriter: c.Writer,
			buff:           &bytes.Buffer{},
		}
		c.Writer = lw

		start := time.Now()
		defer func() {
			end := time.Now()
			logger.CtxInfof(mustGetGinExtractContext(c), "%s|%s|%s|%s|%s|%s|req_body=%s,res_body=%s",
				c.ClientIP(), c.Request.Host, c.Request.Method, c.Request.RequestURI, c.Request.Proto, end.Sub(start),
				request, lw.buff.String(),
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
