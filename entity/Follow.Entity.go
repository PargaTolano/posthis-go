package entity

import (
	"gorm.io/gorm"
)

type Follow struct {
	gorm.Model
	FollowerID uint //User.ID
	FollowedID uint //User.ID
}
