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

	sqlDb, err := db.DB()
	if err != nil {
		return nil, err
	}

	defer sqlDb.Close()

	if err = db.Preload("Likes").First(&post, id).Error; err != nil {
		return nil, err
	}

	return post.Likes, nil
}

func (lm LikeModel) CreateLike(userId, postId uint) (*PostDetailVM, error) {

	user := User{}
	post := Post{}

	db, err := db.ConnectToDb()
	if err != nil {
		return nil, err
	}

	sqlDb, err := db.DB()
	if err != nil {
		return nil, err
	}

	defer sqlDb.Close()

	if err = db.First(&user, userId).Error; err != nil {
		return nil, err
	}

	if err = db.First(&post, postId).Error; err != nil {
		return nil, err
	}

	like := Like{UserID: user.ID, PostID: post.ID}

	//Once it works add everything to the database
	if err = db.Create(&like).Error; err != nil {
		return nil, err
	}

	db.Model(&user).Association("likes").Append(&like)
	db.Model(&post).Association("likes").Append(&like)

	postModel := PostModel(lm)
	model, err := postModel.GetPost(userId, postId)
	if err != nil {
		return nil, err
	}

	return model, nil
}

func (lm LikeModel) DeleteLike(userId, postId uint) (*PostDetailVM, error) {

	like := Like{}
	post := Post{}

	db, err := db.ConnectToDb()
	if err != nil {
		return nil, err
	}

	sqlDb, err := db.DB()
	if err != nil {
		return nil, err
	}

	defer sqlDb.Close()

	if err = db.First(&like, "user_id = ? AND post_id = ?", userId, postId).Error; err != nil {
		return nil, err
	}

	if err = db.First(&post, postId).Error; err != nil {
		return nil, err
	}

	if err = db.Delete(&like).Error; err != nil {
		return nil, err
	}

	postModel := PostModel(lm)
	model, err := postModel.GetPost(userId, postId)
	if err != nil {
		return nil, err
	}

	return model, nil
}
