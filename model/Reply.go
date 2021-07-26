package model

import (
	"gorm.io/gorm"
)

type Reply struct {
	gorm.Model
	content string
	user    User
	post    Post
	media   []PostMedia
}
