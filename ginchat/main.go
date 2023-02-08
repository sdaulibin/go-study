package main

import (
	"ginchat/router"
	"ginchat/utils"

	"github.com/spf13/viper"
)

func main() {
	utils.InitConfig()
	utils.InitMysql()
	utils.InitRedis()
	r := router.Router()
	r.Run(viper.GetString("port.server"))
}
