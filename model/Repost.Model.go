package model

import (
	"posthis/database"
)

type RepostModel struct {
	Model
}

func (RepostModel) GetReposts(id uint) ([]*Repost, error) {

	post := Post{}

	if err := database.DB.Preload("reposts").First(&post, id).Error; err != nil {
		return nil, err
	}

	return post.Reposts, nil
}

func (rm RepostModel) CreateRepost(userId uint, postId uint) (*PostDetailVM, error) {

	user := User{}
	post := Post{}

	if err := database.DB.First(&user, userId).Error; err != nil {
		return nil, err
	}

	if err := database.DB.First(&post, postId).Error; err != nil {
		return nil, err
	}

	repost := Repost{UserID: user.ID, PostID: post.ID}

	//Once it works add everything to the database
	if err := database.DB.Create(&repost).Error; err != nil {
		return nil, err
	}
	database.DB.Model(&user).Association("reposts").Append(&repost)
	database.DB.Model(&post).Association("reposts").Append(&repost)

	postModel := PostModel(rm)
	model, err := postModel.GetPost(userId, postId)
	if err != nil {
		return nil, err
	}

	return model, nil
}

func (rm RepostModel) DeleteRepost(userId, postId uint) (*PostDetailVM, error) {

	repost := Repost{}
	post := Post{}

	if err := database.DB.First(&repost, "user_id = ? AND post_id = ?", userId, postId).Error; err != nil {
		return nil, err
	}

	if err := database.DB.First(&post, postId).Error; err != nil {
		return nil, err
	}

	if err := database.DB.Delete(&repost).Error; err != nil {
		return nil, err
	}

	postModel := PostModel(rm)
	model, err := postModel.GetPost(userId, postId)
	if err != nil {
		return nil, err
	}

	return model, nil
}
