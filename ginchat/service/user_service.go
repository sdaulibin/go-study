package service

import (
	"ginchat/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetUserList
// @Tag 用户列表
// @Success 200 {string} json{"code","message"}
// @Router /user/getUserList [get]
func GetUserList(ctx *gin.Context) {
	data := models.GetUserList()
	ctx.JSON(http.StatusOK, gin.H{
		"data": data,
	})
}
