package viewmodel

type MediaVM struct {
	ID      uint   `json:"id"`
	Path    string `json:"path"`
	Mime    string `json:"mime"`
	IsVideo bool   `json:"isVideo"`
}
