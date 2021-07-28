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

	ProfilePicID *Media    `gorm:"polymorphic:Owner;polymorphicValue:user"`
	CoverPicID   *Media    `gorm:"polymorphic:Owner;polymorphicValue:user"`
	Followers    []*User   `gorm:"many2many:follwing_followers;"`
	Posts        []*Post   `gorm:"polymorphic:Owner;polymorphicValue:user"`
	Likes        []*Like   `gorm:"foreignKey:UserID"`
	Replies      []*Reply  `gorm:"foreignKey:UserID"`
	Reposts      []*Repost `gorm:"foreignKey:UserID"`
}
