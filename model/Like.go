package model

import (
	"gorm.io/gorm"
)

type Like struct {
	gorm.Model
	UserID uint //User.ID
	PostID uint //Post.ID
}
