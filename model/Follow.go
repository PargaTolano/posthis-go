package model

import (
	"gorm.io/gorm"
)

type Follow struct {
	gorm.Model
	followerUser User
	followedUser User
}
