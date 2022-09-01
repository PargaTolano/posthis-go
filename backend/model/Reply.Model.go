package model

import (
	"mime/multipart"
	"posthis/database"
	"posthis/storage"
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
		for _, replymedia := range replies[i].Media {
			mvm := MediaVM{
				ID:      replymedia.ID,
				Path:    replymedia.Url,
				Mime:    replymedia.Mime,
				IsVideo: strings.Contains(replymedia.Mime, "video"),
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

	if len(files) > 0 {
		mediaData, err := storage.UploadMultipleFiles(files)
		if err != nil {
			return nil, err
		}

		for _, data := range mediaData {
			media = append(media, &Media{Name: data.Name, Mime: data.Mime, Url: data.Url})
		}

		database.DB.CreateInBatches(&media, len(media))
	}

	reply := Reply{Content: content, Media: media, UserID: userId, PostID: postId}

	database.DB.Create(&reply)
	database.DB.Model(&user).Association("replies").Append(&reply)
	database.DB.Model(&post).Association("replies").Append(&reply)

	return &reply, nil
}

func (rm ReplyModel) UpdateReply(id uint, content string, deleted []string, files []*multipart.FileHeader) (*ReplyVM, error) {

	user := User{}
	reply := Reply{}

	if err := database.DB.First(&reply, id).Error; err != nil {
		return nil, err
	}

	if err := database.DB.Preload("ProfilePic").Find(&user, reply.UserID).Error; err != nil {
		return nil, err
	}

	if len(deleted) > 0 {
		database.DB.Delete(&Media{}, deleted)
	}

	if len(files) > 0 {
		media := []*Media{}
		mediaData, err := storage.UploadMultipleFiles(files)
		if err != nil {
			return nil, err
		}

		for _, data := range mediaData {
			media = append(media, &Media{Name: data.Name, Mime: data.Mime, Url: data.Url})
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
		model.PublisherProfilePic = user.ProfilePic.Url
	}

	for _, media := range reply.Media {
		mvm := MediaVM{
			ID:      media.ID,
			Path:    media.Url,
			Mime:    media.Mime,
			IsVideo: strings.Contains(media.Mime, "video")}

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
