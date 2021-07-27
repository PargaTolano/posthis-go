package model

import "gorm.io/gorm"

type Media struct {
	gorm.Model
	Mime      string
	Name      string
	OwnerID   int
	OwnerType string
}

func (media *Media) GetPath(scheme, host string) string {
	return scheme + "://" + host + "/path/" + media.Name
}
