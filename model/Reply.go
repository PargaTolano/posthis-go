package model

import (
	"gorm.io/gorm"
)

type Reply struct {
	gorm.Model
	Content string   `gorm:"default:''"`
	UserID  uint     //User.ID
	PostID  uint     //Post.ID
	Media   []*Media `gorm:"polymorphic:Owner;polymorphicValue:reply"`
}
