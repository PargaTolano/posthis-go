package viewmodel

import "time"

type ReplyVM struct {
	ReplyID             uint      `json:"replyID"`
	Content             string    `json:"content"`
	PostID              uint      `json:"postID"`
	PublisherID         uint      `json:"publisherID"`
	PublisherUserName   string    `json:"publisherUserName"`
	PublisherTag        string    `json:"publisherTag"`
	PublisherProfilePic string    `json:"publisherProfilePic"`
	Medias              []MediaVM `json:"medias"`
	Date                time.Time `json:"date"`
}
