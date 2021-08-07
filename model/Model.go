package model

import (
	"posthis/entity"
	"posthis/viewmodel"
)

type Model struct {
	Scheme string
	Host   string
}

type User = entity.User
type Post = entity.Post
type Media = entity.Media
type Like = entity.Like
type Reply = entity.Reply
type Repost = entity.Repost
type Follow = entity.Follow

type PostFeedVM = viewmodel.PostFeedVM
type PostDetailVM = viewmodel.PostDetailVM
type SearchVM = viewmodel.SearchVM
type PostSearchVM = viewmodel.PostSearchVM
type UserSearchVM = viewmodel.UserSearchVM
type UserCreateVM = viewmodel.UserCreateVM
type MediaVM = viewmodel.MediaVM
type UserLoginVm = viewmodel.UserLoginVM
type UserUpdateVM = viewmodel.UserUpdateVM
type UserVM = viewmodel.UserVM
type ReplyVM = viewmodel.ReplyVM
