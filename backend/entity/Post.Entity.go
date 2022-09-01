package entity

import (
	"log"
	"time"

	"gorm.io/gorm"
)

type Post struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Content   string `gorm:"default:''"`
	OwnerID   uint
	OwnerType string
	Media     []*Media  `gorm:"foreignKey:PostOwnerID;constraint:OnDelete:CASCADE;"`
	Likes     []*Like   `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE;"`
	Replies   []*Reply  `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE;"`
	Reposts   []*Repost `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE;"`
}

func (post *Post) BeforeDelete(tx *gorm.DB) (err error) {
	log.Printf("Deleting post: %d...", post.ID)
	return
}
