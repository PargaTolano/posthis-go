package entity

import (
	"log"
	"posthis/storage"
	"time"

	"gorm.io/gorm"
)

type Reply struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Content   string   `gorm:"default:''"`
	UserID    uint     //User.ID
	PostID    uint     //Post.ID
	Media     []*Media `gorm:"foreignKey:ReplyOwnerID;constraint:OnDelete:CASCADE;"`
}

func (reply *Reply) BeforeDelete(tx *gorm.DB) (err error) {
	log.Println(reply.ID, "REPLY Getting Deleted")

	for _, media := range reply.Media {
		// delete from firebase storage
		if err = storage.DeleteFile(media.Name); err != nil {
			return
		}
	}

	return
}
