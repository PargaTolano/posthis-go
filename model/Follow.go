package model

import (
	"gorm.io/gorm"
)

type Follow struct {
	gorm.Model
	FollowerUserID uint //User.ID
	FollowedUserID uint //User.ID
}
