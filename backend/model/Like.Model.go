package model

import (
	"posthis/database"
)

type LikeModel struct {
	Model
}

func (LikeModel) GetLikes(id uint) ([]*Like, error) {

	post := &Post{}

	if err := database.DB.Preload("Likes").First(&post, id).Error; err != nil {
		return nil, err
	}

	return post.Likes, nil
}

func (lm LikeModel) CreateLike(userId, postId uint) (*PostDetailVM, error) {

	user := User{}
	post := Post{}

	if err := database.DB.First(&user, userId).Error; err != nil {
		return nil, err
	}

	if err := database.DB.First(&post, postId).Error; err != nil {
		return nil, err
	}

	like := Like{UserID: user.ID, PostID: post.ID}

	//Once it works add everything to the database
	if err := database.DB.Create(&like).Error; err != nil {
		return nil, err
	}

	database.DB.Model(&user).Association("likes").Append(&like)
	database.DB.Model(&post).Association("likes").Append(&like)

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

	if err := database.DB.First(&like, "user_id = ? AND post_id = ?", userId, postId).Error; err != nil {
		return nil, err
	}

	if err := database.DB.First(&post, postId).Error; err != nil {
		return nil, err
	}

	if err := database.DB.Delete(&like).Error; err != nil {
		return nil, err
	}

	postModel := PostModel(lm)
	model, err := postModel.GetPost(userId, postId)
	if err != nil {
		return nil, err
	}

	return model, nil
}
