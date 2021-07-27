package model

import (
	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	Content   string `gorm:"default:''"`
	OwnerID   uint
	OwnerType string
	Media     []*Media `gorm:"polymorphic:Owner;polymorphicValue:post"`
	Likes     []*Like
	Replies   []*Reply
	Reposts   []*Repost
}
