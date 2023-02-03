package main

import (
	"ginchat/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(mysql.Open("root:binginx@tcp(127.0.0.1:3306)/ginchat?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// 迁移 schema
	//db.AutoMigrate(&models.UserBasic{})
	db.AutoMigrate(&models.Contact{})
	db.AutoMigrate(&models.GroupBasic{})
	db.AutoMigrate(&models.Message{})
	// user := &models.UserBasic{
	// 	Name: "木子木木三",
	// }

	// // Create
	// db.Create(user)

	// d := db.First(user, 1)
	// fmt.Print(d)

	// db.Model(user).Update("Password", "4567")
}
