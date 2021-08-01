package viewmodel

import (
	"time"
)

type PostFeedVM struct {
	ID                  uint
	PublisherID         uint
	PublisherUserName   string
	PublisherTag        string
	PublisherProfilePic string
	ReposterUsername    string
	ReposterTag         string
	Date                time.Time
	Content             string
	Media               []string
	RepostID            uint
}
