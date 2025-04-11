package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func buildRequestMessage(ctx *gin.Context) string {
	var result strings.Builder

	result.WriteString(ctx.ClientIP())
	result.WriteString(" - ")
	result.WriteString(ctx.Request.Method)
	result.WriteString(" ")
	result.WriteString(ctx.Request.RequestURI)
	result.WriteString(" ")

	return result.String()
}

func Logger(l *zap.Logger) func(c *gin.Context) {
	return func(ctx *gin.Context) {
		ctx.Next()
		l.Info(buildRequestMessage(ctx))
	}
}
