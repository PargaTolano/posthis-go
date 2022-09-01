package entity

import (
	"log"
	"posthis/storage"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Tag          string `gorm:"type:varchar(20)"`
	Email        string `gorm:"type:varchar(100)"`
	Username     string `gorm:"type:varchar(20);unique_index"`
	PasswordHash string `gorm:"type:varchar(100)"`

	ProfilePic *Media    `gorm:"foreignKey:ProfilePicOwnerID;constraint:OnDelete:CASCADE;"`
	CoverPic   *Media    `gorm:"foreignKey:CoverPicOwnerID;constraint:OnDelete:CASCADE;"`
	Posts      []*Post   `gorm:"polymorphic:Owner;polymorphicValue:user;constraint:OnDelete:CASCADE;"`
	Followers  []*Follow `gorm:"foreignKey:FollowerID"`
	Following  []*Follow `gorm:"foreignKey:FollowedID"`
	Likes      []*Like   `gorm:"foreignKey:UserID"`
	Replies    []*Reply  `gorm:"foreignKey:UserID"`
	Reposts    []*Repost `gorm:"foreignKey:UserID"`
}

func (user *User) BeforeDelete(tx *gorm.DB) (err error) {
	log.Println(user.ID, "USER Getting Deleted")

	if err = storage.DeleteFile(user.ProfilePic.Name); err != nil {
		return
	}

	if err = storage.DeleteFile(user.CoverPic.Name); err != nil {
		return
	}

	return
}
