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

	db.Joins("JOIN follows ON follows.follower_id = users.id AND follows.follower_id = ?", id).Find(&users)
	if db.Error != nil {
		return nil, err
	}

	return users, nil
}

func (FollowModel) CreateFollow(id, followerId uint) (*Follow, error) {

	user := User{}
	follower := User{}

	db, err := db.ConnectToDb()
	if err != nil {
		return nil, err
	}

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

	return &follow, nil
}

func (FollowModel) DeleteFollow(followerid uint, followedId uint) error {

	followed := User{}
	follower := User{}

	db, err := db.ConnectToDb()
	if err != nil {
		return err
	}

	db.First(&follower, followerid)
	db.First(&followed, followedId)

	db.Model(&follower).Association("Followings").Delete(&followed)
	db.Model(&followed).Association("Followers").Delete(&follower)

	return nil
}
