package model

import (
	"errors"
	"posthis/db"
	"posthis/entity"
	"posthis/viewmodel"
	"strings"
)

type SearchModel struct {
	Model
}

func (sm SearchModel) GetSearch(
	query string,
	searchPost,
	searchUser bool,
	offsetPost,
	limitPost,
	offsetUser,
	limitUser uint) (*SearchVM, error) {

	model := SearchVM{}

	db, err := db.ConnectToDb()
	if err != nil {
		return nil, err
	}

	if !searchPost && !searchUser {
		return nil, errors.New("you have to search at least one of the following fields: { posts, users}")
	}

	if searchPost {

		var (
			ids   []int
			posts []Post
		)

		rows, err := db.Raw("CALL SP_SEARCH_POSTS(?,?,?)", query, offsetPost, limitPost).Rows()
		if err != nil {
			return nil, err
		}

		for rows.Next() {
			var psmodel PostSearchVM
			rows.Scan(
				&psmodel.ID,
				&psmodel.Content,
				&psmodel.PublisherID,
				&psmodel.PublisherUsername,
				&psmodel.PublisherTag,
				&psmodel.PublisherProfilePic,
				&psmodel.Date,
				&psmodel.LikeCount,
				&psmodel.ReplyCount,
				&psmodel.RepostCount)

			psmodel.PublisherProfilePic = entity.GetPath(sm.Scheme, sm.Scheme, psmodel.PublisherProfilePic)

			model.Posts = append(model.Posts, psmodel)
		}
		if !rows.NextResultSet() {
			return nil, rows.Err()
		}
		for rows.Next() {
			var id int
			rows.Scan(&id)
			ids = append(ids, id)
		}

		//relate media to viemodel
		db.Preload("Media").Find(&posts, ids)

		for i := range posts {
			for j := range posts[i].Media {
				model.Posts[i].Media = append(model.Posts[i].Media,
					viewmodel.MediaVM{
						ID:      posts[i].Media[j].ID,
						Path:    posts[i].Media[j].Name,
						Mime:    posts[i].Media[j].Mime,
						IsVideo: strings.Contains(posts[i].Media[j].Mime, "video")})
			}
		}
	}

	if searchUser {
		rows, err := db.Raw("CALL SP_SEARCH_USERS(?,?,?)", query, offsetUser, limitUser).Rows()
		if err != nil {
			return nil, err
		}

		for rows.Next() {
			var umodel UserSearchVM
			rows.Scan(
				&umodel.ID,
				&umodel.Tag,
				&umodel.Username,
				&umodel.ProfilePicPath)

			umodel.ProfilePicPath = entity.GetPath(sm.Scheme, sm.Host, umodel.ProfilePicPath)

			model.Users = append(model.Users, umodel)
		}
	}

	return &model, nil
}
