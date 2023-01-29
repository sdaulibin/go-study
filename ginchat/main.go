package main

import (
	"ginchat/router"
	"ginchat/utils"
)

func main() {
	utils.InitConfig()
	utils.InitMysql()
	r := router.Router()
	r.Run(":8081")
}
