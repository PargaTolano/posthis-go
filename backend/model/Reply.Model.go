package model

import (
	"mime/multipart"
	"posthis/database"
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

	rows, err := database.DB.Raw("CALL SP_GET_POST_REPLIES(?)", id).Rows()
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

	database.DB.Preload("Media").Clauses(clause.OrderBy{
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

	if err := database.DB.First(&user, userId).Error; err != nil {
		return nil, err
	}

	if err := database.DB.First(&post, postId).Error; err != nil {
		return nil, err
	}

	err := utils.UploadMultipleFiles(files, &media)
	if err != nil {
		return nil, err
	}
	database.DB.CreateInBatches(&media, len(media))

	reply := Reply{Content: content, Media: media, UserID: userId, PostID: postId}

	database.DB.Create(&reply)
	database.DB.Model(&user).Association("replies").Append(&reply)
	database.DB.Model(&post).Association("replies").Append(&reply)

	return &reply, nil
}

func (rm ReplyModel) UpdateReply(id uint, content string, deleted []string, files []*multipart.FileHeader) (*ReplyVM, error) {

	user := User{}
	reply := Reply{}
	deletedMedia := []*Media{}

	if err := database.DB.First(&reply, id).Error; err != nil {
		return nil, err
	}

	if err := database.DB.Preload("ProfilePic").Find(&user, reply.UserID).Error; err != nil {
		return nil, err
	}

	if len(deleted) > 0 {
		database.DB.Find(&deletedMedia, deleted)
		for _, dm := range deletedMedia {
			utils.DeleteStaticFile(dm.Name)
		}
		database.DB.Delete(&deletedMedia)
	}

	if len(files) >= 1 {
		media := []*Media{}
		err := utils.UploadMultipleFiles(files, &media)
		if err != nil {
			return nil, err
		}

		database.DB.CreateInBatches(&media, len(media))
		database.DB.Model(&reply).Association("Media").Append(media)
	}

	if content != "" {
		reply.Content = content
	}

	database.DB.Save(&reply)

	model := ReplyVM{
		ReplyID:           reply.ID,
		Content:           reply.Content,
		PostID:            reply.PostID,
		PublisherID:       reply.UserID,
		PublisherUserName: user.Username,
		Date:              reply.CreatedAt}

	if user.ProfilePic != nil {
		model.PublisherProfilePic = user.ProfilePic.GetPath(rm.Scheme, rm.Host)
	}

	for i := range reply.Media {

		mvm := MediaVM{
			ID:      reply.Media[i].ID,
			Path:    reply.Media[i].GetPath(rm.Scheme, rm.Host),
			Mime:    reply.Media[i].Mime,
			IsVideo: strings.Contains(reply.Media[i].Mime, "video")}

		model.Medias = append(model.Medias, mvm)
	}

	return &model, nil
}

func (ReplyModel) DeleteReply(id uint) error {

	if err := database.DB.Delete(&Reply{}, id).Error; err != nil {
		return err
	}

	return nil
}
