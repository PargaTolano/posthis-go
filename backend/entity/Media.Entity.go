package entity

import (
	"time"
)

type Media struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
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

func GetPath(scheme, host, name string) string {

	if len(name) == 0 {
		return ""
	}

	if len(scheme) == 0 {
		scheme = "http"
	}

	if len(host) == 0 {
		host = "localhost:4000"
	}

	return scheme + "://" + host + "/static/" + name
}
