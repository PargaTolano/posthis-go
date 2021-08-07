package model

import (
	"mime/multipart"
	"posthis/db"
	"posthis/entity"
	"posthis/utils"
	"strings"

	"gorm.io/gorm/clause"
)

type ReplyModel struct {
	Model
}

func (rm ReplyModel) GetReplies(id uint) ([]ReplyVM, error) {

	replies := []Reply{}
	models := []ReplyVM{}
	replyIds := []uint{}

	db, err := db.ConnectToDb()
	if err != nil {
		return nil, err
	}

	rows, err := db.Raw("CALL SP_GET_POST_REPLIES(?)", id).Rows()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		model := ReplyVM{}
		rows.Scan(
			&model.ReplyID,
			&model.Content,
			&model.PostID,
			&model.PublisherID,
			&model.PublisherUserName,
			&model.PublisherTag,
			&model.PublisherProfilePic,
			&model.Date)

		model.PublisherProfilePic = entity.GetPath(rm.Scheme, rm.Host, model.PublisherProfilePic)

		models = append(models, model)
	}
	if !rows.NextResultSet() {
		return nil, rows.Err()
	}

	//Obtener id de posts para obtener su media
	for rows.Next() {
		var id uint
		rows.Scan(&id)
		replyIds = append(replyIds, id)
	}

	db.Preload("Media").Clauses(clause.OrderBy{
		Expression: clause.Expr{SQL: "FIELD(id,?)", Vars: []interface{}{replyIds}, WithoutParentheses: true},
	}).Find(&replies, replyIds)

	for i := range replies {
		for j := range replies[i].Media {
			mvm := MediaVM{
				ID:      replies[i].Media[j].ID,
				Path:    replies[i].Media[j].GetPath(rm.Scheme, rm.Host),
				Mime:    replies[i].Media[j].Mime,
				IsVideo: strings.Contains(replies[i].Media[j].Mime, "video"),
			}

			models[i].Medias = append(models[i].Medias, mvm)
		}
	}

	return models, nil
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
