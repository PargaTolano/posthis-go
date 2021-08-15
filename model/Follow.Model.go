package model

import (
	"posthis/database"
	"posthis/entity"
)

type FollowModel struct {
	Model
}

func (fm FollowModel) GetFollows(id, viewerId uint) ([]FollowUserVM, error) {

	models := []FollowUserVM{}

	rows, err := database.DB.Raw("CALL SP_GET_FOLLOWERS(?,?)", id, viewerId).Rows()
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

	rows, err := database.DB.Raw("CALL SP_GET_FOLLOWING(?,?)", id, viewerId).Rows()
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

	database.DB.First(&user, id)
	database.DB.First(&follower, followerId)

	if database.DB.Error != nil {
		return nil, database.DB.Error
	}

	follow := Follow{FollowerID: follower.ID, FollowedID: user.ID}
	if err := database.DB.FirstOrCreate(&follow).Error; err != nil {
		return nil, err
	}
	database.DB.Model(&user).Association("followers").Append(&follow)
	database.DB.Model(&follower).Association("followings").Append(&follow)

	model, err := userModel.GetUser(id, followerId)
	if err != nil {
		return nil, err
	}

	return model, nil
}

func (fm FollowModel) DeleteFollow(followerId uint, followedId uint) (*UserVM, error) {

	userModel := UserModel(fm)
	follow := Follow{}

	if err := database.DB.First(&follow, "follower_id = ? AND followed_id = ?", followerId, followedId).Error; err != nil {
		return nil, err
	}

	database.DB.Delete(&follow)

	model, err := userModel.GetUser(followedId, followerId)
	if err != nil {
		return nil, err
	}

	return model, nil
}
