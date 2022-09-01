package model

import (
	"errors"
	"mime/multipart"
	"net/http"
	"posthis/auth"
	"posthis/database"
	"posthis/storage"
	"posthis/utils"
)

type UserModel struct {
	Model
}

func (UserModel) GetUsers() ([]User, error) {
	var (
		users []User
	)

	database.DB.Preload("ProfilePic").Preload("CoverPic").Find(&users)
	if database.DB.Error != nil {
		return nil, database.DB.Error
	}

	return users, nil
}

func (um UserModel) GetUser(id, viewerId uint) (*UserVM, error) {

	model := UserVM{}

	row := database.DB.Raw("CALL SP_GET_PROFILE(?,?)", id, viewerId).Row()

	row.Scan(
		&model.ID,
		&model.Username,
		&model.Tag,
		&model.Email,
		&model.ProfilePicPath,
		&model.CoverPicPath,
		&model.FollowerCount,
		&model.FollowingCount,
		&model.IsFollowed)

	return &model, nil
}

func (UserModel) CreateUser(model UserCreateVM) (*User, error) {
	hash, err := utils.HashPassword(model.Password)

	if err != nil {
		return nil, err
	}

	user := User{Tag: model.Tag, Email: model.Email, Username: model.Username, PasswordHash: hash}
	database.DB.Create(&user)
	return &user, nil
}

func (um UserModel) UpdateUser(id, viewerId uint, model UserUpdateVM, pfpFiles []*multipart.FileHeader, coverFiles []*multipart.FileHeader) (*UserVM, error) {

	user := User{}

	if err := database.DB.Preload("ProfilePic").Preload("CoverPic").First(&user, id).Error; err != nil {
		return nil, err
	}

	if model.Tag != "" {
		user.Tag = model.Tag
	}

	if model.Email != "" {
		user.Email = model.Email
	}

	if model.Username != "" {
		user.Username = model.Username
	}

	if len(pfpFiles) == 1 {
		media := []*Media{}

		if user.ProfilePic != nil {
			database.DB.Delete(&user.ProfilePic)
		}

		mediaData, err := storage.UploadMultipleFiles(pfpFiles)
		if err != nil {
			return nil, err
		}

		for _, data := range mediaData {
			media = append(media, &Media{Name: data.Name, Mime: data.Mime, Url: data.Mime})
		}

		database.DB.CreateInBatches(&media, len(media))
		database.DB.Model(&user).Association("ProfilePic").Append(&media)
	}

	if len(coverFiles) == 1 {
		media := []*Media{}

		if user.CoverPic != nil {
			database.DB.Delete(&user.CoverPic)
		}

		mediaData, err := storage.UploadMultipleFiles(coverFiles)
		if err != nil {
			return nil, err
		}

		for _, data := range mediaData {
			media = append(media, &Media{Name: data.Name, Mime: data.Mime, Url: data.Mime})
		}
		database.DB.CreateInBatches(&media, len(media))
		database.DB.Model(&user).Association("CoverPic").Append(&media)
	}

	database.DB.Save(&user)

	umodel, err := um.GetUser(id, viewerId)
	if err != nil {
		return nil, err
	}

	return umodel, nil
}

func (UserModel) DeleteUser(id uint) error {

	database.DB.Delete(&User{}, id)
	if database.DB.Error != nil {
		return database.DB.Error
	}

	return nil
}

func (um UserModel) ValidatePassword(id uint, password string) (bool, error) {
	user := User{}

	if err := database.DB.First(&user, id).Error; err != nil {
		return false, err
	}

	return utils.CheckPasswordHash(password, user.PasswordHash), nil
}

//return auth token info and error if there is one
func (um UserModel) Login(model UserLoginVm) (interface{}, error) {

	user := User{Username: model.Username}

	if err := database.DB.Preload("ProfilePic").Find(&user).Error; err != nil {
		return nil, err
	}

	if !utils.CheckPasswordHash(model.Password, user.PasswordHash) {
		return nil, errors.New("wrong password")
	}

	ts, err := auth.CreateToken(uint64(user.ID))
	if err != nil {
		return nil, err
	}

	err = auth.CreateAuth(uint64(user.ID), ts)
	if err != nil {
		return nil, err
	}

	ppp := ""
	if user.ProfilePic != nil {
		ppp = user.ProfilePic.Url
	}

	data := map[string]interface{}{
		"id":             user.ID,
		"username":       user.Username,
		"token":          ts.AccessToken,
		"profilePicPath": ppp,
	}

	return data, nil
}

//return auth token info and error if there is one
func (UserModel) Logout(r *http.Request) error {

	au, err := auth.ExtractTokenMetadata(r)
	if err != nil {
		return err
	}

	deleted, err := auth.DeleteAuth(au.AccessUuid)
	if err != nil || deleted == 0 {
		return err
	}

	return nil
}
