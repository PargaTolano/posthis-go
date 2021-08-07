package viewmodel

type UserVM struct {
	ID             uint   `json:"id"`
	Username       string `json:"username"`
	Tag            string `json:"tag"`
	Email          string `json:"email"`
	ProfilePicPath string `json:"profilePicPath"`
	CoverPicPath   string `json:"coverPicPath"`
	FollowerCount  int    `json:"followerCount"`
	FollowingCount int    `json:"followingCount"`
	IsFollowed     bool   `json:"isFollowed"`
}
