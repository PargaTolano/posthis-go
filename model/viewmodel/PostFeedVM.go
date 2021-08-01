package viewmodel

import (
	"time"
)

type PostFeedVM struct {
	ID                  uint      `json:"postID"`
	PublisherID         uint      `json:"publisherID"`
	PublisherUserName   string    `json:"publisherUserName"`
	PublisherTag        string    `json:"publisherTag"`
	PublisherProfilePic string    `json:"publisherProfilePic"`
	ReposterUsername    string    `json:"reposterUserName"`
	ReposterTag         string    `json:"reposterTag"`
	Date                time.Time `json:"date"`
	Content             string    `json:"content"`
	Media               []MediaVM `json:"medias"`
	RepostID            uint      `json:"isRepost"`
}
