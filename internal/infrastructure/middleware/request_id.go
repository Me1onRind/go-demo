package middleware

import (
	"github.com/Me1onRind/go-demo/internal/infrastructure/logger"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

func SetRequestId() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestId := uuid.NewV4().String()
		setGinExtractContext(c, logger.WithFields(mustGetGinExtractContext(c), logrus.Fields{
			logger.RequestIdKey: requestId,
		}))
		c.Next()
	}
}
