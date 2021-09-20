package entity

import (
	"time"
)

type Reply struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Content   string   `gorm:"default:''"`
	UserID    uint     //User.ID
	PostID    uint     //Post.ID
	Media     []*Media `gorm:"polymorphic:Owner;polymorphicValue:reply"`
}
