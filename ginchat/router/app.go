package router

import (
	"ginchat/docs"
	"ginchat/service"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Router() *gin.Engine {
	r := gin.Default()
	//swagger
	docs.SwaggerInfo.BasePath = ""
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	//static resource
	r.Static("/asset", "asset/")
	r.LoadHTMLGlob("views/**/*")

	//index
	r.GET("/", service.GetIndex)
	r.GET("/index", service.GetIndex)
	r.GET("/user/register", service.Register)
	r.GET("/toChat", service.ToChat)
	r.GET("/user/searchFriends", service.SearchFriends)

	//user
	r.GET("/user/getUserList", service.GetUserList)
	r.POST("/user/createUser", service.CreateUser)
	r.GET("/user/deleteUser", service.DeleteUser)
	r.POST("/user/updateUser", service.UpdateUser)
	r.POST("/user/login", service.GetUser)

	//message
	r.GET("/user/sendMsg", service.SendMsg)
	r.GET("/user/sendUserMsg", service.SendUserMsg)
	return r
}
