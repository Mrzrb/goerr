package demo

import (
	"time"

	"github.com/gin-gonic/gin"
)

// @Logger
func TestPackage(ctx *gin.Context, t time.Time) int {
	return 1
}
