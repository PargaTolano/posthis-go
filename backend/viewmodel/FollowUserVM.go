package viewmodel

type FollowUserVM struct {
	ID             uint   `json:"id"`
	Username       string `json:"username"`
	Tag            string `json:"tag"`
	ProfilePicPath string `json:"profilePicPath"`
	IsFollowed     bool   `json:"isFollowed"`
}
