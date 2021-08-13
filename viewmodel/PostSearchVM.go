package viewmodel

import (
	"time"
)

type PostSearchVM struct {
	ID                  uint      `json:"postID"`
	Content             string    `json:"content"`
	Media               []MediaVM `json:"medias"`
	PublisherID         uint      `json:"publisherID"`
	PublisherUsername   string    `json:"publisherUserName"`
	PublisherTag        string    `json:"publisherTag"`
	PublisherProfilePic string    `json:"publisherProfilePic"`
	Date                time.Time `json:"date"`
	LikeCount           int       `json:"likeCount"`
	ReplyCount          int       `json:"replyCount"`
	RepostCount         int       `json:"repostCount"`
	IsLiked             int       `json:"isLiked"`
	IsReposted          int       `json:"isReposted"`
}
