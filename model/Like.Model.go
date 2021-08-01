package model

import (
	"posthis/db"
)

type LikeModel struct {
	Model
}

func (LikeModel) GetLikes(id uint) ([]*Like, error) {

	post := &Post{}

	db, err := db.ConnectToDb()
	if err != nil {
		return nil, err
	}

	if err = db.Preload("Likes").First(&post, id).Error; err != nil {
		return nil, err
	}

	return post.Likes, nil
}

func (LikeModel) CreateLike(userId, postId uint) (*Like, error) {

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

	like := Like{UserID: user.ID, PostID: post.ID}

	//Once it works add everything to the database
	if err = db.FirstOrCreate(&like).Error; err != nil {
		return nil, err
	}

	db.Model(&user).Association("likes").Append(&like)
	db.Model(&post).Association("likes").Append(&like)

	return &like, nil
}

func (LikeModel) DeleteLike(id uint) error {

	db, err := db.ConnectToDb()
	if err != nil {
		return err
	}

	if err = db.Delete(&Like{}, id).Error; err != nil {
		return err
	}

	return nil
}
