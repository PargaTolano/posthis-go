package model

import (
	"posthis/db"
	"posthis/entity"
)

type FollowModel struct {
	Model
}

func (fm FollowModel) GetFollows(id, viewerId uint) ([]FollowUserVM, error) {

	models := []FollowUserVM{}

	db, err := db.ConnectToDb()
	if err != nil {
		return nil, err
	}

	sqlDb, err := db.DB()
	if err != nil {
		return nil, err
	}

	defer sqlDb.Close()

	rows, err := db.Raw("CALL SP_GET_FOLLOWERS(?,?)", id, viewerId).Rows()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		model := FollowUserVM{}
		rows.Scan(
			&model.ID,
			&model.Username,
			&model.Tag,
			&model.ProfilePicPath,
			&model.IsFollowed,
		)

		model.ProfilePicPath = entity.GetPath(fm.Scheme, fm.Host, model.ProfilePicPath)

		models = append(models, model)
	}
	return models, nil
}

func (fm FollowModel) GetFollowing(id, viewerId uint) ([]FollowUserVM, error) {
	models := []FollowUserVM{}

	db, err := db.ConnectToDb()
	if err != nil {
		return nil, err
	}

	sqlDb, err := db.DB()
	if err != nil {
		return nil, err
	}

	defer sqlDb.Close()

	rows, err := db.Raw("CALL SP_GET_FOLLOWING(?,?)", id, viewerId).Rows()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		model := FollowUserVM{}
		rows.Scan(
			&model.ID,
			&model.Username,
			&model.Tag,
			&model.ProfilePicPath,
			&model.IsFollowed,
		)

		model.ProfilePicPath = entity.GetPath(fm.Scheme, fm.Host, model.ProfilePicPath)

		models = append(models, model)
	}
	return models, nil
}

func (fm FollowModel) CreateFollow(id, followerId uint) (*UserVM, error) {

	userModel := UserModel(fm)
	user := User{}
	follower := User{}

	db, err := db.ConnectToDb()
	if err != nil {
		return nil, err
	}

	sqlDb, err := db.DB()
	if err != nil {
		return nil, err
	}

	defer sqlDb.Close()

	db.First(&user, id)
	db.First(&follower, followerId)

	if db.Error != nil {
		return nil, db.Error
	}

	follow := Follow{FollowerID: follower.ID, FollowedID: user.ID}
	if err = db.FirstOrCreate(&follow).Error; err != nil {
		return nil, err
	}
	db.Model(&user).Association("followers").Append(&follow)
	db.Model(&follower).Association("followings").Append(&follow)

	model, err := userModel.GetUser(id, followerId)
	if err != nil {
		return nil, err
	}

	return model, nil
}

func (fm FollowModel) DeleteFollow(followerId uint, followedId uint) (*UserVM, error) {

	userModel := UserModel(fm)
	follow := Follow{}

	db, err := db.ConnectToDb()
	if err != nil {
		return nil, err
	}

	sqlDb, err := db.DB()
	if err != nil {
		return nil, err
	}

	defer sqlDb.Close()

	if err = db.First(&follow, "follower_id = ? AND followed_id = ?", followerId, followedId).Error; err != nil {
		return nil, err
	}

	db.Delete(&follow)

	model, err := userModel.GetUser(followedId, followerId)
	if err != nil {
		return nil, err
	}

	return model, nil
}
