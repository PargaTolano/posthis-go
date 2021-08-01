package db

import "posthis/entity"

type User = entity.User
type Post = entity.Post
type Media = entity.Media
type Like = entity.Like
type Reply = entity.Reply
type Repost = entity.Repost
type Follow = entity.Follow

func InitDB() {
	db, err := ConnectToDb()
	if err != nil {
		panic("Error: " + err.Error())
	}

	//Add models to db
	db.AutoMigrate(&User{}, &Post{}, &Media{}, &Like{}, &Reply{}, &Repost{}, &Follow{})
}
