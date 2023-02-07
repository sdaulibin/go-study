package models

import (
	"ginchat/utils"

	"gorm.io/gorm"
)

type Contact struct {
	gorm.Model
	OwnerId  uint //谁的关系信息
	TargetId uint //对应的谁
	Type     int  //1-好友 2-群 3-xxx
}

func (table *Contact) TableName() string {
	return "contact"
}

func SearchFriends(userId uint) []UserBasic {
	contacts := make([]Contact, 0)
	users := make([]UserBasic, 0)
	utils.DB.Where("owner_id = ? and type = 1", userId).Find(&contacts)
	for _, v := range contacts {
		user := FindUserByID(uint(v.TargetId))
		users = append(users, user)
	}
	return users
}

func AddFriend(userId uint, targetName string) int {
	user := UserBasic{}
	if targetName != "" {
		user = FindUserByName(targetName)
		if user.Salt != "" {
			contact := Contact{}
			contact.OwnerId = userId
			contact.TargetId = user.ID
			contact.Type = 1
			utils.DB.Create(&contact)
			return 0
		}
		return -1
	}
	return -1
}
