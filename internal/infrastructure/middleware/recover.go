package middleware

import "github.com/gin-gonic/gin"

func Recover() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := recover(); err != nil {
		}
	}
}
