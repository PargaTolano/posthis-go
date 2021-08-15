package model

import (
	"errors"
	"mime/multipart"
	"net/http"
	"posthis/auth"
	"posthis/database"
	"posthis/entity"
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

	model.ProfilePicPath = entity.GetPath(um.Scheme, um.Host, model.ProfilePicPath)
	model.CoverPicPath = entity.GetPath(um.Scheme, um.Host, model.CoverPicPath)

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
			utils.DeleteStaticFile(user.ProfilePic.Name)
			database.DB.Delete(&user.ProfilePic)
		}

		utils.UploadMultipleFiles(pfpFiles, &media)
		database.DB.CreateInBatches(&media, len(media))
		database.DB.Model(&user).Association("ProfilePic").Append(&media)
	}

	if len(coverFiles) == 1 {
		media := []*Media{}

		if user.CoverPic != nil {
			utils.DeleteStaticFile(user.CoverPic.Name)
			database.DB.Delete(&user.CoverPic)
		}

		utils.UploadMultipleFiles(coverFiles, &media)
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

	user := User{}

	if err := database.DB.Preload("ProfilePic").Where("username = ?", model.Username).Find(&user).Error; err != nil {
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
		ppp = user.ProfilePic.GetPath(um.Scheme, um.Host)
	}

	data := map[string]interface{}{
		"id":             user.ID,
		"username":       user.Username,
		"profilePicPath": ppp,
		"token":          ts.AccessToken,
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
