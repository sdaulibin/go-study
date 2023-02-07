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
			tx := utils.DB.Begin()
			contact1 := Contact{}
			contact1.OwnerId = userId
			contact1.TargetId = user.ID
			contact1.Type = 1
			utils.DB.Create(&contact1)
			contact2 := Contact{}
			contact2.OwnerId = user.ID
			contact2.TargetId = userId
			contact2.Type = 1
			utils.DB.Create(&contact2)
			tx.Commit()
			return 0
		}
		return -1
	}
	return -1
}
