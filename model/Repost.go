package model

import (
	"gorm.io/gorm"
)

type Repost struct {
	gorm.Model
	user User
	post Post
}
