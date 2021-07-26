package model

import (
	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	content string
	user    User
	media   []PostMedia
	likes   []Like
	replies []Reply
	reposts []Repost
}
