package v1

import (
	"github.com/gin-gonic/gin"
)

func errorResponse(ctx *gin.Context, code int, msg string) error {
	ctx.JSON(code, gin.H{"error": msg})
	return nil
}
