package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Tag          string `gorm:"type:varchar(20)"`
	Email        string `gorm:"type:varchar(100)"`
	Username     string `gorm:"type:varchar(20);unique_index"`
	PasswordHash string `gorm:"type:varchar(100)"`

	ProfilePic *Media    `gorm:"polymorphic:Owner;polymorphicValue:user"`
	CoverPic   *Media    `gorm:"polymorphic:Owner;polymorphicValue:user"`
	Posts      []*Post   `gorm:"polymorphic:Owner;polymorphicValue:user"`
	Followers  []*Follow `gorm:"foreignKey:FollowedID"`
	Following  []*Follow `gorm:"foreignKey:FollowerID"`
	Likes      []*Like   `gorm:"foreignKey:UserID"`
	Replies    []*Reply  `gorm:"foreignKey:UserID"`
	Reposts    []*Repost `gorm:"foreignKey:UserID"`
}
