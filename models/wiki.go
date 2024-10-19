package models

import "gorm.io/gorm"

type WikiPage struct {
	gorm.Model
	Path  string `gorm:"unique"`
	Title string `gorm:"index"`
	Body  string
}

type Note struct {
	gorm.Model
	Body string
}
