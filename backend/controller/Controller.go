package controller

import (
	"posthis/model"
	"posthis/viewmodel"
)

type Model = model.Model
type UserModel = model.UserModel
type PostModel = model.PostModel
type MediaModel = model.MediaModel
type LikeModel = model.LikeModel
type ReplyModel = model.ReplyModel
type RepostModel = model.RepostModel
type FollowModel = model.FollowModel
type SearchModel = model.SearchModel

type SearchVM = viewmodel.SearchVM
type PostSearchVM = viewmodel.PostSearchVM
type UserSearchVM = viewmodel.UserSearchVM
type UserCreateVM = viewmodel.UserCreateVM
type UserUpdateVM = viewmodel.UserUpdateVM
type UserLoginVM = viewmodel.UserLoginVM
type SuccesVM = viewmodel.SuccesVM
type ReplyUpdateVM = viewmodel.ReplyUpdateVM
