package models

import (
	"fmt"
	"ginchat/utils"

	"gorm.io/gorm"
)

type Community struct {
	gorm.Model
	Name    string
	OwnerId uint
	Icon    string
	Cute    string
	Desc    string
}

func (table *Community) TableName() string {
	return "community"
}

func CreateCommunity(community Community) (int, string) {
	if len(community.Name) == 0 {
		return -1, "群名称不能为空"
	}
	if community.OwnerId == 0 {
		return -1, "请先登录"
	}
	// user := FindUserByID(community.OwnerId)
	// if user.Salt == "" {
	// 	return -1, "应户名错误"
	// }
	if err := utils.DB.Create(&community).Error; err != nil {
		fmt.Println(err)
		return -1, "建群失败"
	}
	return 0, "建群成功"
}

func GetCommunities(ownerId uint) ([]*Community, string) {
	contacts := make([]Contact, 0)
	ids := make([]uint, 0)
	utils.DB.Where("owner_id = ? and type = 2", ownerId).Find(&contacts)
	for _, v := range contacts {
		ids = append(ids, v.TargetId)
	}
	data := make([]*Community, 0)
	utils.DB.Where("id in ?", ids).Find(&data)
	return data, "查询完成"
}

func JoinGroup(userId uint, communityId string) (int, string) {
	contact := Contact{}
	contact.OwnerId = userId
	// contact.TargetId = commId
	contact.Type = 2
	community := Community{}

	utils.DB.Where("id = ? or name = ?", communityId, communityId).Find(&community)
	if community.Name == "" {
		return -1, "没找到群"
	}
	temp := FindContact(userId, community.ID, 2)
	if !temp.CreatedAt.IsZero() {
		return -1, "已加过此群"
	} else {
		contact.TargetId = community.ID
		utils.DB.Create(&contact)
		return 0, "已群成功"
	}
}
