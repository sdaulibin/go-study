package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetIndex
// @Tag 首页
// @Success 200 {string} welcome
// @Router /index [get]
func GetIndex(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"meaasge": "welcome hello !!!!",
	})
}
