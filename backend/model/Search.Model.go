package model

import (
	"errors"
	"posthis/database"
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
	viewerId,
	offsetUser,
	limitUser uint) (*SearchVM, error) {

	model := SearchVM{}
	model.Users = make([]UserSearchVM, 0)
	model.Posts = make([]PostSearchVM, 0)

	if !searchPost && !searchUser {
		return nil, errors.New("you have to search at least one of the following fields: { posts, users}")
	}

	if searchPost {

		var (
			ids   = make([]int, 0)
			posts = make([]Post, 0)
		)

		rows, err := database.DB.Raw("CALL SP_SEARCH_POSTS(?,?,?,?)", query, viewerId, offsetPost, limitPost).Rows()
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
				&psmodel.RepostCount,
				&psmodel.IsLiked,
				&psmodel.IsReposted)

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
		database.DB.Preload("Media").Find(&posts, ids)

		for i := range model.Posts {
			for j := range posts[i].Media {
				model.Posts[i].Media = append(model.Posts[i].Media,
					viewmodel.MediaVM{
						ID:      posts[i].Media[j].ID,
						Path:    posts[i].Media[j].GetPath(sm.Scheme, sm.Host),
						Mime:    posts[i].Media[j].Mime,
						IsVideo: strings.Contains(posts[i].Media[j].Mime, "video")})
			}
		}
	}

	if searchUser {

		rows, err := database.DB.Raw("CALL SP_SEARCH_USERS(?,?,?)", query, offsetUser, limitUser).Rows()
		if err != nil {
			return nil, err
		}

		for rows.Next() {
			var umodel UserSearchVM
			rows.Scan(
				&umodel.ID,
				&umodel.Tag,
				&umodel.Username,
				&umodel.ProfilePicPath,
				&umodel.FollowerCount,
				&umodel.FollowingCount)

			umodel.ProfilePicPath = entity.GetPath(sm.Scheme, sm.Host, umodel.ProfilePicPath)

			model.Users = append(model.Users, umodel)
		}
	}

	return &model, nil
}
