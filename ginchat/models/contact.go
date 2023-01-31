package models

import "gorm.io/gorm"

type Contact struct {
	gorm.Model
	OwnerId  int64
	TargetId int64
	Type     int
}

func (table *Contact) TableName() string {
	return "contact"
}
