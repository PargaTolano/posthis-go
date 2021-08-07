package viewmodel

import "time"

type PostSearchVM struct {
	ID                  uint
	Content             string
	Media               []MediaVM
	PublisherID         uint
	PublisherUsername   string
	PublisherTag        string
	PublisherProfilePic string
	Date                time.Time
	LikeCount           int
	ReplyCount          int
	RepostCount         int
}
