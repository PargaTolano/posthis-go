package entity

import (
	"time"
)

type Media struct {
	ID                uint `gorm:"primarykey"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
	Mime              string
	Name              string
	Url               string
	ProfilePicOwnerID *uint `gorm:"uniqueIndex:idx_string_owner_id_other_owner_id"`
	CoverPicOwnerID   *uint `gorm:"uniqueIndex:idx_string_owner_id_other_owner_id"`
	PostOwnerID       *uint `gorm:"uniqueIndex:idx_string_owner_id_other_owner_id"`
	ReplyOwnerID      *uint `gorm:"uniqueIndex:idx_string_owner_id_other_owner_id"`
}
