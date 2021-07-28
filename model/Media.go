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

	if len(scheme) == 0 {
		scheme = "http"
	}

	if len(host) == 0 {
		host = "localhost:4000"
	}

	return scheme + "://" + host + "/static/" + media.Name
}
