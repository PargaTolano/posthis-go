package entity

import (
	"time"
)

type Follow struct {
	ID         uint `gorm:"primarykey"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	FollowerID uint //User.ID
	FollowedID uint //User.ID
}
