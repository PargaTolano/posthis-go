package model

import (
	"posthis/db"
)

type RepostModel struct {
	Model
}

func (RepostModel) GetReposts(id uint) ([]*Repost, error) {

	post := Post{}

	db, err := db.ConnectToDb()
	if err != nil {
		return nil, err
	}

	if err = db.Preload("reposts").First(&post, id).Error; err != nil {
		return nil, err
	}

	return post.Reposts, nil
}

func (RepostModel) CreateRepost(userId uint, postId uint) (*Repost, error) {

	user := User{}
	post := Post{}

	db, err := db.ConnectToDb()
	if err != nil {
		return nil, err
	}

	if err = db.First(&user, userId).Error; err != nil {
		return nil, err
	}

	if err = db.First(&post, postId).Error; err != nil {
		return nil, err
	}

	repost := Repost{UserID: user.ID, PostID: post.ID}

	//Once it works add everything to the database
	if err = db.FirstOrCreate(&repost).Error; err != nil {
		return nil, err
	}
	db.Model(&user).Association("reposts").Append(&repost)
	db.Model(&post).Association("reposts").Append(&repost)

	return &repost, nil
}

func (RepostModel) DeleteRepost(id uint) error {

	repost := Repost{}

	db, err := db.ConnectToDb()
	if err != nil {
		return err
	}

	if err = db.Delete(&repost, id).Error; err != nil {
		return err
	}

	return nil
}
