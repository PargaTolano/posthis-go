package model

import (
	"errors"
	"mime/multipart"
	"net/http"
	"posthis/auth"
	"posthis/db"
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

	db, err := db.ConnectToDb()
	if err != nil {
		return nil, err
	}

	sqlDb, err := db.DB()
	if err != nil {
		return nil, err
	}

	defer sqlDb.Close()

	db.Preload("ProfilePic").Preload("CoverPic").Find(&users)
	if db.Error != nil {
		return nil, db.Error
	}

	return users, nil
}

func (um UserModel) GetUser(id, viewerId uint) (*UserVM, error) {

	model := UserVM{}

	db, err := db.ConnectToDb()
	if err != nil {
		return nil, err
	}

	sqlDb, err := db.DB()
	if err != nil {
		return nil, err
	}

	defer sqlDb.Close()

	row := db.Raw("CALL SP_GET_PROFILE(?,?)", id, viewerId).Row()

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

	db, err := db.ConnectToDb()
	if err != nil {
		return nil, err
	}

	sqlDb, err := db.DB()
	if err != nil {
		return nil, err
	}

	defer sqlDb.Close()

	user := User{Tag: model.Tag, Email: model.Email, Username: model.Username, PasswordHash: hash}
	db.Create(&user)
	return &user, nil
}

func (um UserModel) UpdateUser(id, viewerId uint, model UserUpdateVM, pfpFiles []*multipart.FileHeader, coverFiles []*multipart.FileHeader) (*UserVM, error) {

	user := User{}

	db, err := db.ConnectToDb()
	if err != nil {
		return nil, err
	}

	sqlDb, err := db.DB()
	if err != nil {
		return nil, err
	}

	defer sqlDb.Close()

	if err = db.Preload("ProfilePic").Preload("CoverPic").First(&user, id).Error; err != nil {
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
			db.Delete(&user.ProfilePic)
		}

		utils.UploadMultipleFiles(pfpFiles, &media)
		db.CreateInBatches(&media, len(media))
		db.Model(&user).Association("ProfilePic").Append(&media)
	}

	if len(coverFiles) == 1 {
		media := []*Media{}

		if user.CoverPic != nil {
			utils.DeleteStaticFile(user.CoverPic.Name)
			db.Delete(&user.CoverPic)
		}

		utils.UploadMultipleFiles(coverFiles, &media)
		db.CreateInBatches(&media, len(media))
		db.Model(&user).Association("CoverPic").Append(&media)
	}

	db.Save(&user)

	umodel, err := um.GetUser(id, viewerId)
	if err != nil {
		return nil, err
	}

	return umodel, nil
}

func (UserModel) DeleteUser(id uint) error {
	db, err := db.ConnectToDb()
	if err != nil {
		return err
	}

	sqlDb, err := db.DB()
	if err != nil {
		return err
	}

	defer sqlDb.Close()

	db.Delete(&User{}, uint(id))
	if db.Error != nil {
		return db.Error
	}

	return nil
}

func (um UserModel) ValidatePassword(id uint, password string) (bool, error) {
	user := User{}

	db, err := db.ConnectToDb()
	if err != nil {
		return false, err
	}

	sqlDb, err := db.DB()
	if err != nil {
		return false, err
	}

	defer sqlDb.Close()

	if err = db.First(&user, id).Error; err != nil {
		return false, err
	}

	return utils.CheckPasswordHash(password, user.PasswordHash), nil
}

//return auth token info and error if there is one
func (um UserModel) Login(model UserLoginVm) (interface{}, error) {

	user := User{}

	db, err := db.ConnectToDb()
	if err != nil {
		return nil, err
	}

	sqlDb, err := db.DB()
	if err != nil {
		return nil, err
	}

	defer sqlDb.Close()

	if err := db.Preload("ProfilePic").Where("username = ?", model.Username).Find(&user).Error; err != nil {
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
