package model

import (
	"mime/multipart"
	"posthis/database"
	"posthis/storage"
	"strings"

	"gorm.io/gorm/clause"
)

type PostModel struct {
	Model
}

func (PostModel) GetPosts() ([]Post, error) {

	posts := []Post{}

	database.DB.Preload("Media").Find(posts)

	return posts, nil
}

func (pm PostModel) GetPost(userId, id uint) (*PostDetailVM, error) {
	post := Post{}
	model := PostDetailVM{}

	query := `CALL SP_GET_POST_DETAIL(?, ?)`

	err := database.DB.Raw(query, id, userId).Row().Scan(
		&model.ID,
		&model.PublisherID,
		&model.PublisherUserName,
		&model.PublisherTag,
		&model.PublisherProfilePic,
		&model.ReposterUsername,
		&model.ReposterTag,
		&model.Date,
		&model.Content,
		&model.RepostID,
		&model.LikeCount,
		&model.ReplyCount,
		&model.RepostCount,
		&model.IsLiked,
		&model.IsReposted)

	if err != nil {
		return nil, err
	}

	database.DB.Preload("Media").First(&post, id)

	for _, media := range post.Media {
		mvm := MediaVM{
			ID:      media.ID,
			Path:    media.Url,
			Mime:    media.Mime,
			IsVideo: strings.Contains("video", media.Mime)}
		model.Media = append(model.Media, mvm)
	}

	return &model, nil
}

func (PostModel) CreatePost(ownerId uint, content string, files []*multipart.FileHeader) (*Post, error) {

	var (
		poster User
		post   Post
		media  []*Media
	)

	mediaData, err := storage.UploadMultipleFiles(files)
	if err != nil {
		return nil, err
	}

	// create database media from data
	for _, data := range mediaData {
		media = append(media, &Media{Name: data.Name, Mime: data.Mime, Url: data.Url})
	}

	if err := database.DB.CreateInBatches(&media, len(media)).Error; err != nil {
		println("failed to create media", err.Error())
		return nil, err
	}

	post = Post{Content: content, Media: media}

	if err := database.DB.First(&poster, ownerId).Error; err != nil {
		println("failed to find poster, poster ID: ", ownerId, err.Error())
		return nil, err
	}

	if err := database.DB.Create(&post).Error; err != nil {
		println("failed to create post", err.Error())
		return nil, err
	}

	if err := database.DB.Model(&poster).Association("Posts").Append(&post); err != nil {
		println("failed to associate post to poster", err.Error())
		return nil, err
	}

	return &post, nil
}

func (pm PostModel) UpdatePost(userId, id uint, content string, deleted []string, files []*multipart.FileHeader) (*PostDetailVM, error) {

	var (
		post         Post
		deletedMedia []*Media
	)

	database.DB.First(&post, id)

	database.DB.Find(&deletedMedia, deleted)
	database.DB.Delete(&deletedMedia)

	if len(files) > 0 {
		var media []*Media
		mediaData, err := storage.UploadMultipleFiles(files)
		if err != nil {
			return nil, err
		}

		for _, data := range mediaData {
			media = append(media, &Media{Name: data.Name, Mime: data.Mime, Url: data.Url})
		}

		database.DB.CreateInBatches(&media, len(media))
		database.DB.Model(&post).Association("Media").Append(media)
	}

	if content != "" {
		post.Content = content
	}

	database.DB.Save(&post)

	model, err := pm.GetPost(userId, post.ID)
	if err != nil {
		return nil, err
	}

	return model, nil
}

// TODO MANUAL STORAGE DELETE ON ALL MEDIA DELETING FUNCTIONS
func (PostModel) DeletePost(id uint) error {
	var post Post
	database.DB.Model(&Post{}).Preload("Media").First(&post, id)

	for _, media := range post.Media {
		if err := storage.DeleteFile(media.Name); err != nil {
			return err
		}
	}

	database.DB.Delete(&post)
	if err := database.DB.Error; err != nil {
		return err
	}

	return nil
}

func (pm PostModel) GetFeed(id, offset, limit uint) ([]PostFeedVM, error) {

	models := []PostFeedVM{}
	postIds := []uint{}
	posts := []Post{}

	query := `CALL SP_GET_FEED(?,?,?)`

	rows, err := database.DB.Raw(query, id, offset, limit).Rows()
	if err != nil {
		return nil, err
	}
	//posts
	for rows.Next() {
		model := PostFeedVM{}
		rows.Scan(
			&model.ID,
			&model.PublisherID,
			&model.PublisherUserName,
			&model.PublisherTag,
			&model.PublisherProfilePic,
			&model.ReposterUsername,
			&model.ReposterTag,
			&model.Date,
			&model.Content,
			&model.RepostID,
			&model.LikeCount,
			&model.ReplyCount,
			&model.RepostCount,
			&model.IsLiked,
			&model.IsReposted)

		models = append(models, model)
	}

	if len(models) == 0 {
		return make([]PostFeedVM, 0), nil
	}

	if !rows.NextResultSet() {
		return nil, rows.Err()
	}

	//Obtener id de posts para obtener su media
	for rows.Next() {
		var id uint
		rows.Scan(&id)
		postIds = append(postIds, id)
	}

	database.DB.Preload("Media").Clauses(clause.OrderBy{
		Expression: clause.Expr{SQL: "FIELD(id,?)", Vars: []interface{}{postIds}, WithoutParentheses: true},
	}).Find(&posts, postIds)

	for i := range posts {
		for _, postmedia := range posts[i].Media {
			mvm := MediaVM{
				ID:      postmedia.ID,
				Path:    postmedia.Url,
				Mime:    postmedia.Mime,
				IsVideo: strings.Contains(postmedia.Mime, "video"),
			}

			models[i].Media = append(models[i].Media, mvm)
		}
	}

	return models, nil
}

func (pm PostModel) GetUserFeed(id, userId, offset, limit uint) ([]PostFeedVM, error) {

	models := []PostFeedVM{}
	postIds := []uint{}
	posts := []Post{}

	query := `CALL SP_GET_USER_FEED(?,?,?,?)`

	rows, err := database.DB.Raw(query, id, userId, offset, limit).Rows()
	if err != nil {
		return nil, err
	}
	//posts
	for rows.Next() {
		model := PostFeedVM{}
		rows.Scan(
			&model.ID,
			&model.PublisherID,
			&model.PublisherUserName,
			&model.PublisherTag,
			&model.PublisherProfilePic,
			&model.ReposterUsername,
			&model.ReposterTag,
			&model.Date,
			&model.Content,
			&model.RepostID,
			&model.LikeCount,
			&model.ReplyCount,
			&model.RepostCount,
			&model.IsLiked,
			&model.IsReposted)
		models = append(models, model)
	}
	if !rows.NextResultSet() {
		return nil, rows.Err()
	}

	//Obtener id de posts para obtener su media
	for rows.Next() {
		var id uint
		rows.Scan(&id)
		postIds = append(postIds, id)
	}

	database.DB.Preload("Media").Clauses(clause.OrderBy{
		Expression: clause.Expr{SQL: "FIELD(id,?)", Vars: []interface{}{postIds}, WithoutParentheses: true},
	}).Find(&posts, postIds)

	for i := range posts {
		for _, postmedia := range posts[i].Media {
			mvm := MediaVM{
				ID:      postmedia.ID,
				Path:    postmedia.Url,
				Mime:    postmedia.Mime,
				IsVideo: strings.Contains(postmedia.Mime, "video"),
			}

			models[i].Media = append(models[i].Media, mvm)
		}
	}

	return models, nil
}
