package repository

import (
	"github.com/4otis/library_api_2025/internal/models"
	"gorm.io/gorm"
)

type AuthorRepository struct {
	db *gorm.DB
}

func NewAuthorRepository(db *gorm.DB) *AuthorRepository {
	return &AuthorRepository{db: db}
}

func (ar AuthorRepository) Create(author *models.Author) error {
	return ar.db.Create(author).Error
}

func (ar AuthorRepository) Read(id uint) (author *models.Author, err error) {
	err = ar.db.Preload("Books").First(&author, id).Error
	return author, err
}

func (ar AuthorRepository) ReadAll() (authors []*models.Author, err error) {
	err = ar.db.Preload("Books").Find(&authors).Error
	return authors, err
}

func (ar AuthorRepository) Update(id uint, newAuthor *models.Author) error {
	return ar.db.Transaction(func(tx *gorm.DB) error {
		var author models.Author
		if err := tx.First(&author, id).Error; err != nil {
			return err
		}

		if err := tx.Model(&author).Updates(newAuthor).Error; err != nil {
			return err
		}

		if newAuthor.Books != nil {
			err := tx.Model(&author).Association("Books").Replace(newAuthor.Books)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func (ar AuthorRepository) Delete(id uint) error {
	return ar.db.Select("Books").Delete(&models.Author{Model: gorm.Model{ID: id}}).Error
}
