package service

import (
	"fmt"
	"ginchat/models"
	"net/http"
	"strconv"

	"github.com/asaskevich/govalidator"
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
// @param id query string false "用户ID"
// @Success 200 {string}  json{"code","message"}
// @Router /user/deleteUser [get]
func DeleteUser(ctx *gin.Context) {
	user := models.UserBasic{}
	idStr := ctx.Query("id")
	if idStr == "" {
		ctx.JSON(-1, gin.H{
			"message": "用户名id不能为空",
		})
		return
	}
	id, _ := strconv.Atoi(idStr)
	user.ID = uint(id)
	models.DeleteUser(user)
	ctx.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("删除用户id>%s成功！", idStr),
	})
}

// UpdateUser
// @Tags 用户模块
// @Summary 用户修改
// @param id formData string false "用户ID"
// @param name formData string false "用户名"
// @param password formData string false "密码"
// @param phone formData string false "手机号"
// @param email formData string false "邮箱"
// @Success 200 {string}  json{"code","message"}
// @Router /user/updateUser [post]
func UpdateUser(ctx *gin.Context) {
	user := models.UserBasic{}
	id, _ := strconv.Atoi(ctx.PostForm("id"))
	user.ID = uint(id)
	name := ctx.PostForm("name")
	user.Name = name
	password := ctx.PostForm("password")
	user.Password = password
	phone := ctx.PostForm("phone")
	user.Phone = phone
	email := ctx.PostForm("email")
	user.Email = email
	_, err := govalidator.ValidateStruct(user)
	if err != nil {
		fmt.Println("err:", err)
		ctx.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("修改用户>%s失败: %s！", name, err),
		})
		return
	}
	models.UpdateUser(user)
	ctx.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("修改用户%s成功！", name),
	})
}
