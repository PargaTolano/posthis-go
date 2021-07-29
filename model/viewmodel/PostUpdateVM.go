package viewmodel

type PostUpdateVM struct {
	Content string `json:"content"`
	Deleted []int  `json:"deleted"`
}
