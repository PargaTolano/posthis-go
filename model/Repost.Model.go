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

	sqlDb, err := db.DB()
	if err != nil {
		return nil, err
	}

	defer sqlDb.Close()

	if err = db.Preload("reposts").First(&post, id).Error; err != nil {
		return nil, err
	}

	return post.Reposts, nil
}

func (rm RepostModel) CreateRepost(userId uint, postId uint) (*PostDetailVM, error) {

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

	repost := Repost{UserID: user.ID, PostID: post.ID}

	//Once it works add everything to the database
	if err = db.Create(&repost).Error; err != nil {
		return nil, err
	}
	db.Model(&user).Association("reposts").Append(&repost)
	db.Model(&post).Association("reposts").Append(&repost)

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

	db, err := db.ConnectToDb()
	if err != nil {
		return nil, err
	}

	sqlDb, err := db.DB()
	if err != nil {
		return nil, err
	}

	defer sqlDb.Close()

	if err = db.First(&repost, "user_id = ? AND post_id = ?", userId, postId).Error; err != nil {
		return nil, err
	}

	if err = db.First(&post, postId).Error; err != nil {
		return nil, err
	}

	if err = db.Delete(&repost).Error; err != nil {
		return nil, err
	}

	postModel := PostModel(rm)
	model, err := postModel.GetPost(userId, postId)
	if err != nil {
		return nil, err
	}

	return model, nil
}
