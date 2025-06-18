package models

import "gorm.io/gorm"

type Book struct {
	gorm.Model
	Title   string    `json:"title"`
	Pages   int       `json:"pages"`
	Authors []*Author `json:"authors" gorm:"many2many:books_authors;"`
}
