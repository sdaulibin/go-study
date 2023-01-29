package service

import (
	"fmt"
	"ginchat/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetUserList
// @Tags 用户模块
// @Summary 用户列表
// @Success 200 {string}  json{"code","message"}
// @Router /user/getUserList [get]
func GetUserList(ctx *gin.Context) {
	data := models.GetUserList()
	ctx.JSON(http.StatusOK, gin.H{
		"data": data,
	})
}

// CreateUser
// @Tags 用户模块
// @Summary 用户新增
// @param name query string false "用户名"
// @param password query string false "密码"
// @param repassword query string false "确认密码"
// @Success 200 {string}  json{"code","message"}
// @Router /user/createUser [get]
func CreateUser(ctx *gin.Context) {
	user := models.UserBasic{}
	user.Name = ctx.Query("name")
	password := ctx.Query("password")
	repassword := ctx.Query("repassword")
	if password != repassword {
		ctx.JSON(-1, gin.H{
			"message": "两次密码不一致",
		})
		return
	}
	user.Password = repassword
	models.CreateUser(user)
	ctx.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("新增用户%s成功！", user.Name),
	})
}

// DeleteUser
// @Tags 用户模块
// @Summary 用户删除
// @param name query string false "用户名"
// @Success 200 {string}  json{"code","message"}
// @Router /user/deleteUser [get]
func DeleteUser(ctx *gin.Context) {
	user := models.UserBasic{}
	name := ctx.Query("name")
	if name == "" {
		ctx.JSON(-1, gin.H{
			"message": "用户名不能为空",
		})
		return
	}
	user.Name = name
	models.DeleteUser(user)
	ctx.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("删除用户%s成功！", name),
	})
}
