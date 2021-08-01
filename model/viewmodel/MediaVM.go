package viewmodel

type MediaVM struct {
	ID      uint   `json:"id,omitempty"`
	Path    string `json:"path,omitempty"`
	Mime    string `json:"mime,omitempty"`
	IsVideo bool   `json:"is,omitempty"`
}
