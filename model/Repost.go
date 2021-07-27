package model

import (
	"gorm.io/gorm"
)

type Repost struct {
	gorm.Model
	UserID uint //User.ID
	PostID uint //Post.ID
}
