package model

import "gorm.io/gorm"

type Media struct {
	gorm.Model
	mime string
	name string
}
