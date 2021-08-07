package viewmodel

type ReplyUpdateVM struct {
	Content string `json:"content"`
	Deleted []int  `json:"deleted"`
}
