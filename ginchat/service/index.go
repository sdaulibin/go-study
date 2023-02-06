package service

import (
	"fmt"
	"ginchat/models"
	"html/template"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetIndex
// @Tags 首页
// @Success 200 {string} welcome
// @Router /index [get]
func GetIndex(ctx *gin.Context) {
	// ctx.JSON(http.StatusOK, gin.H{
	// 	"meaasge": "welcome hello !!!!",
	// })
	index, err := template.ParseFiles("index.html", "views/chat/head.html")
	if err != nil {
		panic(err)
	}
	fmt.Println(index)
	index.Execute(ctx.Writer, "index")
}

func Register(ctx *gin.Context) {
	register, err := template.ParseFiles("views/user/register.html")
	if err != nil {
		panic(err)
	}
	register.Execute(ctx.Writer, "register")
}

func ToChat(ctx *gin.Context) {
	chat, err := template.ParseFiles("views/chat/index.html",
		"views/chat/head.html",
		"views/chat/concat.html",
		"views/chat/group.html",
		"views/chat/tabmenu.html",
		"views/chat/main.html",
		"views/chat/createcom.html",
		"views/chat/userinfo.html",
		"views/chat/foot.html",
		"views/chat/profile.html")
	if err != nil {
		panic(err)
	}
	user := models.UserBasic{}
	userId, _ := strconv.Atoi(ctx.Query("userId"))
	token := ctx.Query("token")
	user.ID = uint(userId)
	user.Identity = token
	chat.Execute(ctx.Writer, user)
}

func Chat(ctx *gin.Context) {
	models.Chat(ctx.Writer, ctx.Request)
}
