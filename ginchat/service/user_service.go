package service

import (
	"fmt"
	"ginchat/models"
	"ginchat/utils"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// GetUserList
// @Tags 用户模块
// @Summary 用户列表
// @Success 200 {string}  json{"code","message"}
// @Router /user/getUserList [get]
func GetUserList(ctx *gin.Context) {
	data := models.GetUserList()
	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "查新列表成功！",
		"data":    data,
	})
}

// GetUser
// @Tags 用户模块
// @Summary 获取用户
// @param name formData string false "用户名"
// @param password formData string false "密码"
// @Success 200 {string}  json{"code","message"}
// @Router /user/getUser [post]
func GetUser(ctx *gin.Context) {
	name := ctx.Request.FormValue("name")
	password := ctx.Request.FormValue("password")
	user := models.FindUserByName(name)
	if user.Name == "" {
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  fmt.Sprintf("用户>%s 错误！！", name),
		})
		return
	}
	flag := utils.ValidatePwd(password, user.Salt, user.Password)
	if !flag {
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "密码错误！",
		})
		return
	}
	models.UpdateUserToken(user)
	data := models.FindUserByNameAndPwd(name, utils.MakePassword(password, user.Salt))
	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  fmt.Sprintf("获取用户%s成功！", user.Name),
		"data": data,
	})
}

// CreateUser
// @Tags 用户模块
// @Summary 用户新增
// @param name formData string false "用户名"
// @param password formData string false "密码"
// @param repassword formData string false "确认密码"
// @Success 200 {string}  json{"code","message"}
// @Router /user/createUser [post]
func CreateUser(ctx *gin.Context) {
	user := models.UserBasic{}
	user.Name = ctx.Request.FormValue("name")
	password := ctx.Request.FormValue("password")
	repassword := ctx.Request.FormValue("repassword")
	if user.Name == "" || password == "" || repassword == "" {
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "用户名或密码不能为空",
		})
		return
	}
	if password != repassword {
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "两次密码不一致",
		})
		return
	}
	salt := fmt.Sprintf("%06d", rand.Int31())
	user.Salt = salt
	user.Password = utils.MakePassword(password, salt)
	data := models.FindUserByName(user.Name)
	if data.Name != "" {
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  fmt.Sprintf("用户>%s<已注册！", user.Name),
		})
		return
	}
	user.Identity = utils.MD5Encode(fmt.Sprintf("%d", time.Now().Unix()))
	models.CreateUser(user)
	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  fmt.Sprintf("新增用户>%s<成功！", user.Name),
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
		ctx.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"message": "用户名id不能为空",
		})
		return
	}
	id, _ := strconv.Atoi(idStr)
	user.ID = uint(id)
	models.DeleteUser(user)
	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
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
			"code":    -1,
			"message": fmt.Sprintf("修改用户>%s失败: %s！", name, err),
		})
		return
	}
	models.UpdateUser(user)
	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": fmt.Sprintf("修改用户%s成功！", name),
	})
}

var upGrade = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func SendMsg(ctx *gin.Context) {
	ws, err := upGrade.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer func(ws *websocket.Conn) {
		err = ws.Close()
		if err != nil {
			fmt.Println(err)
			return
		}
	}(ws)
	msgHandler(ws, ctx)
}

func msgHandler(ws *websocket.Conn, ctx *gin.Context) {
	fmt.Println("ctx: ", ctx)
	msg, err := utils.Subscribe(ctx, utils.PUBLISH_KEY)
	if err != nil {
		fmt.Println(err)
		return
	}
	tm := time.Now().Format("2006-01-02 15:04:05")
	message := fmt.Sprintf("[ws][%s]:%s", tm, msg)
	fmt.Println("send message: ", message)
	err = ws.WriteMessage(1, []byte(message))
	if err != nil {
		fmt.Println(err)
		return
	}
}

func SendUserMsg(ctx *gin.Context) {
	models.Chat(ctx.Writer, ctx.Request)
}

func SearchFriends(ctx *gin.Context) {
	userId, _ := strconv.Atoi(ctx.Request.FormValue("userId"))
	users := models.SearchFriends(uint(userId))
	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  fmt.Sprintf("获取好友列表成功！>%d<", userId),
		"data": users,
		"Rows": users,
	})
}
