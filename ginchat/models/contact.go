package models

import (
	"fmt"
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

func FindContact(ownerId, targetId uint, ctype int) Contact {
	contact := Contact{}
	utils.DB.Where("owner_id = ? and target_id = ? and type = ?", ownerId, targetId, ctype).Find(&contact)
	return contact
}

func AddFriend(userId uint, targetName string) (int, string) {
	user := UserBasic{}
	if targetName != "" {
		user = FindUserByName(targetName)
		if user.Salt != "" {
			if userId == user.ID {
				return -1, fmt.Sprintf("不能加自己！>%s<", targetName)
			}
			contact := FindContact(userId, user.ID, 1)
			if contact.ID != 0 {
				return -1, fmt.Sprintf("已添加！>%s<", targetName)
			}
			tx := utils.DB.Begin()
			defer func() {
				if r := recover(); r != nil {
					tx.Rollback()
				}
			}()
			contact1 := Contact{}
			contact1.OwnerId = userId
			contact1.TargetId = user.ID
			contact1.Type = 1
			if err := utils.DB.Create(&contact1).Error; err != nil {
				tx.Rollback()
				return -1, fmt.Sprintf("添加好友失败！>%s<", targetName)
			}
			contact2 := Contact{}
			contact2.OwnerId = user.ID
			contact2.TargetId = userId
			contact2.Type = 1
			utils.DB.Create(&contact2)
			if err := utils.DB.Create(&contact2).Error; err != nil {
				tx.Rollback()
				return -1, fmt.Sprintf("添加好友失败！>%s<", targetName)
			}
			tx.Commit()
			return 0, fmt.Sprintf("添加好友成功！>%s<", targetName)
		}
		return -1, fmt.Sprintf("添加好友失败！>%s<", targetName)
	}
	return -1, fmt.Sprintf("添加好友失败！>%s<", targetName)
}
