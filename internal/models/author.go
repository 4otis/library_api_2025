package models

import "gorm.io/gorm"

type Author struct {
	gorm.Model
	Name  string  `json:"name"`
	Books []*Book `json:"books" gorm:"many2many:books_authors;"`
}
