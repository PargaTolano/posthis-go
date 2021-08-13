package model

import (
	"posthis/db"
)

type FollowModel struct {
	Model
}

func (FollowModel) GetFollows(id uint) ([]User, error) {

	users := []User{}

	db, err := db.ConnectToDb()
	if err != nil {
		return nil, err
	}

	sqlDb, err := db.DB()
	if err != nil {
		return nil, err
	}

	defer sqlDb.Close()

	db.Joins("JOIN follows ON follows.follower_id = users.id AND follows.followed_id = ?", id).Find(&users)
	if db.Error != nil {
		return nil, err
	}

	return users, nil
}

func (FollowModel) GetFollowing(id uint) ([]User, error) {
	users := []User{}

	db, err := db.ConnectToDb()
	if err != nil {
		return nil, err
	}

	sqlDb, err := db.DB()
	if err != nil {
		return nil, err
	}

	defer sqlDb.Close()

	db.Joins("JOIN follows ON follows.follower_id = users.id AND follows.follower_id = ?", id).Find(&users)
	if db.Error != nil {
		return nil, err
	}

	return users, nil
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
