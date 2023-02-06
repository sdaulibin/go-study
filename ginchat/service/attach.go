package service

import (
	"fmt"
	"ginchat/utils"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func Upload(ctx *gin.Context) {
	w := ctx.Writer
	req := ctx.Request
	srcFile, header, err := req.FormFile("file")
	if err != nil {
		utils.RespFail(w, err.Error())
	}
	siffix := ".png"
	ofilName := header.Filename
	temp := strings.Split(ofilName, ".")
	if len(temp) > 1 {
		siffix = "." + temp[len(temp)-1]
	}

	fileName := fmt.Sprintf("%d%04d%s", time.Now().Unix(), rand.Int31(), siffix)
	dstFile, err := os.Create("./asset/upload/" + fileName)
	if err != nil {
		utils.RespFail(w, err.Error())
	}
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		utils.RespFail(w, err.Error())
	}
	url := "./asset/upload/" + fileName
	utils.RespOk(w, url, "send file success")
}
