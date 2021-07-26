package model

import (
	"gorm.io/gorm"
)

type Like struct {
	gorm.Model
	user User
	post Post
}
