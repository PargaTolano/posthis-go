package model

import (
	"mime/multipart"
	"posthis/db"
	"posthis/utils"
)

type ReplyModel struct {
	Model
}

func (ReplyModel) GetReplies(id uint) ([]*Reply, error) {

	post := Post{}

	db, err := db.ConnectToDb()
	if err != nil {
		return nil, err
	}

	if err = db.Preload("Replies").First(&post, id).Error; err != nil {
		return nil, err
	}

	return post.Replies, nil
}

func (ReplyModel) CreateReply(userId, postId uint, content string, files []*multipart.FileHeader) (*Reply, error) {

	user := User{}
	post := Post{}
	media := []*Media{}

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

	err = utils.UploadMultipleFiles(files, &media)
	if err != nil {
		return nil, err
	}
	db.CreateInBatches(&media, len(media))

	reply := Reply{Content: content, Media: media, UserID: userId, PostID: postId}

	//Once it works add everything to the database
	db.Create(&reply)
	db.Model(&user).Association("replies").Append(&reply)
	db.Model(&post).Association("replies").Append(&reply)

	return &reply, nil
}

func (ReplyModel) UpdateReply(id uint, content string, deleted []string, files []*multipart.FileHeader) (*Reply, error) {

	reply := Reply{}
	deletedMedia := []*Media{}

	db, err := db.ConnectToDb()
	if err != nil {
		return nil, err
	}

	if err := db.First(&reply, id).Error; err != nil {
		return nil, err
	}

	if err = db.Find(&reply, id).Error; err != nil {
		return nil, err
	}

	db.Find(&deletedMedia, deleted)
	for _, dm := range deletedMedia {
		utils.DeleteStaticFile(dm.Name)
	}
	db.Delete(&deletedMedia)

	if len(files) >= 1 {
		media := []*Media{}
		err = utils.UploadMultipleFiles(files, &media)
		if err != nil {
			return nil, err
		}

		db.CreateInBatches(&media, len(media))
		db.Model(&reply).Association("Media").Append(media)
	}

	if content != "" {
		reply.Content = content
	}

	db.Save(&reply)

	return &reply, nil
}

func (ReplyModel) DeleteReply(id uint) error {

	db, err := db.ConnectToDb()
	if err != nil {
		return err
	}

	if err = db.Delete(&Reply{}, id).Error; err != nil {
		return err
	}

	return nil
}
