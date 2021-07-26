package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	tag          string
	email        string
	username     string
	passwordHash string

	profilePic UserMedia
	coverPic   UserMedia
	followers  []User
	following  []User
	posts      []Post
	likes      []Like
	replies    []Reply
	reposts    []Repost
}
