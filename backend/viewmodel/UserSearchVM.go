package viewmodel

type UserSearchVM struct {
	ID             uint   `json:"id"`
	Tag            string `json:"tag"`
	Username       string `json:"username"`
	ProfilePicPath string `json:"profilePicPath"`
	FollowerCount  uint   `json:"followerCount"`
	FollowingCount uint   `json:"followingCount"`
}
