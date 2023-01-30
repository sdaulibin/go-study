package models

import "gorm.io/gorm"

type Contact struct {
	gorm.Model
	OwnerId  uint
	TargetId uint
	Type     int
}

func (table *Contact) TableName() string {
	return "contact"
}
