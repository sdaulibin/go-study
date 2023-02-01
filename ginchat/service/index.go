package service

import (
	"fmt"
	"html/template"

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
	index, err := template.ParseFiles("index.html")
	if err != nil {
		panic(err)
	}
	fmt.Println(index)
	index.Execute(ctx.Writer, "index")
}
