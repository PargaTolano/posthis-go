package viewmodel

type SearchVM struct {
	Users []UserSearchVM `json:"users"`
	Posts []PostSearchVM `json:"posts"`
}
