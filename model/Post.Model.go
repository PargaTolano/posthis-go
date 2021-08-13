package model

import (
	"mime/multipart"
	"posthis/db"
	"posthis/entity"
	"posthis/utils"
	"strings"

	"gorm.io/gorm/clause"
)

type PostModel struct {
	Model
}

func (PostModel) GetPosts() ([]Post, error) {

	posts := []Post{}

	db, err := db.ConnectToDb()
	if err != nil {
		return nil, err
	}

	sqlDb, err := db.DB()
	if err != nil {
		return nil, err
	}

	defer sqlDb.Close()

	db.Preload("Media").Find(posts)

	return posts, nil
}

func (pm PostModel) GetPost(userId, id uint) (*PostDetailVM, error) {

	post := Post{}
	model := PostDetailVM{}

	db, err := db.ConnectToDb()
	if err != nil {
		return nil, err
	}

	sqlDb, err := db.DB()
	if err != nil {
		return nil, err
	}

	defer sqlDb.Close()

	query := `CALL SP_GET_POST_DETAIL(?, ?)`

	err = db.Raw(query, id, userId).Row().Scan(
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

	model.PublisherProfilePic = entity.GetPath(pm.Scheme, pm.Host, model.PublisherProfilePic)

	db.Preload("Media").First(&post, id)

	for i := range post.Media {

		mvm := MediaVM{
			ID:      post.Media[i].ID,
			Path:    post.Media[i].GetPath(pm.Scheme, pm.Host),
			Mime:    post.Media[i].Mime,
			IsVideo: strings.Contains("video", post.Media[i].Mime)}
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

	db, err := db.ConnectToDb()
	if err != nil {
		return nil, err
	}

	sqlDb, err := db.DB()
	if err != nil {
		return nil, err
	}

	defer sqlDb.Close()

	db.First(&poster, ownerId)

	err = utils.UploadMultipleFiles(files, &media)
	if err != nil {
		return nil, err
	}

	db.CreateInBatches(&media, len(media))

	post = Post{Content: content, Media: media}

	db.Create(&post)
	db.Model(&poster).Association("Posts").Append(&post)

	return &post, nil
}

func (pm PostModel) UpdatePost(userId, id uint, content string, deleted []string, files []*multipart.FileHeader) (*PostDetailVM, error) {

	var (
		post         Post
		media        []*Media
		deletedMedia []*Media
	)

	db, err := db.ConnectToDb()
	if err != nil {
		return nil, err
	}

	sqlDb, err := db.DB()
	if err != nil {
		return nil, err
	}

	defer sqlDb.Close()

	db.First(&post, id)

	db.Find(&deletedMedia, deleted)
	for _, dm := range deletedMedia {
		utils.DeleteStaticFile(dm.Name)
	}
	db.Delete(&deletedMedia)

	if len(files) >= 1 {
		err = utils.UploadMultipleFiles(files, &media)
		if err != nil {
			return nil, err
		}

		db.CreateInBatches(&media, len(media))
		db.Model(&post).Association("Media").Append(media)
	}

	if content != "" {
		post.Content = content
	}

	db.Save(&post)

	model, err := pm.GetPost(userId, post.ID)
	if err != nil {
		return nil, err
	}

	return model, nil
}

func (PostModel) DeletePost(id uint) error {

	db, err := db.ConnectToDb()
	if err != nil {
		return err
	}

	sqlDb, err := db.DB()
	if err != nil {
		return err
	}

	defer sqlDb.Close()

	db.Delete(&Post{}, id)
	if err := db.Error; err != nil {
		return err
	}

	return nil
}

func (pm PostModel) GetFeed(id, offset, limit uint) ([]PostFeedVM, error) {

	models := []PostFeedVM{}
	postIds := []uint{}
	posts := []Post{}

	db, err := db.ConnectToDb()
	if err != nil {
		return nil, err
	}

	sqlDb, err := db.DB()
	if err != nil {
		return nil, err
	}

	defer sqlDb.Close()

	query := `CALL SP_GET_FEED(?,?,?)`

	rows, err := db.Raw(query, id, offset, limit).Rows()
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

		model.PublisherProfilePic = entity.GetPath(pm.Scheme, pm.Host, model.PublisherProfilePic)

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

	db.Preload("Media").Clauses(clause.OrderBy{
		Expression: clause.Expr{SQL: "FIELD(id,?)", Vars: []interface{}{postIds}, WithoutParentheses: true},
	}).Find(&posts, postIds)

	for i := range posts {
		for j := range posts[i].Media {
			mvm := MediaVM{
				ID:      posts[i].Media[j].ID,
				Path:    posts[i].Media[j].GetPath(pm.Scheme, pm.Host),
				Mime:    posts[i].Media[j].Mime,
				IsVideo: strings.Contains(posts[i].Media[j].Mime, "video"),
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

	db, err := db.ConnectToDb()
	if err != nil {
		return nil, err
	}

	sqlDb, err := db.DB()
	if err != nil {
		return nil, err
	}

	defer sqlDb.Close()

	query := `CALL SP_GET_USER_FEED(?,?,?,?)`

	rows, err := db.Raw(query, id, userId, offset, limit).Rows()
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

		model.PublisherProfilePic = entity.GetPath(pm.Scheme, pm.Host, model.PublisherProfilePic)

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

	db.Preload("Media").Clauses(clause.OrderBy{
		Expression: clause.Expr{SQL: "FIELD(id,?)", Vars: []interface{}{postIds}, WithoutParentheses: true},
	}).Find(&posts, postIds)

	for i := range posts {
		for j := range posts[i].Media {
			mvm := MediaVM{
				ID:      posts[i].Media[j].ID,
				Path:    posts[i].Media[j].GetPath(pm.Scheme, pm.Host),
				Mime:    posts[i].Media[j].Mime,
				IsVideo: strings.Contains(posts[i].Media[j].Mime, "video"),
			}

			models[i].Media = append(models[i].Media, mvm)
		}
	}

	return models, nil
}
