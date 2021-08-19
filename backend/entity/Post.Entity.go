package entity

import (
	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	Content   string `gorm:"default:''"`
	OwnerID   uint
	OwnerType string
	Media     []*Media  `gorm:"polymorphic:Owner;polymorphicValue:post"`
	Likes     []*Like   `gorm:"foreignKey:PostID"`
	Replies   []*Reply  `gorm:"foreignKey:PostID"`
	Reposts   []*Repost `gorm:"foreignKey:PostID"`
}
