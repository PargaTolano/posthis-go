package entity

import (
	"time"
)

type Repost struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    uint //User.ID
	PostID    uint //Post.ID
}
