package repository

import (
	"github.com/4otis/library_api_2025/internal/models"
	"gorm.io/gorm"
)

type BookRepository struct {
	db *gorm.DB
}

func NewBookRepository(db *gorm.DB) *BookRepository {
	return &BookRepository{db: db}
}

func (br BookRepository) Create(book *models.Book) error {
	return br.db.Create(book).Error
}

func (br BookRepository) Read(id uint) (book *models.Book, err error) {
	err = br.db.Preload("Authors").First(&book, id).Error
	return book, err
}

func (br BookRepository) ReadAll() (books []*models.Book, err error) {
	err = br.db.Preload("Authors").Find(&books).Error
	return books, err
}

func (br BookRepository) Update(id uint, newBook *models.Book) error {
	return br.db.Transaction(func(tx *gorm.DB) error {
		var book models.Book
		if err := tx.First(&book, id).Error; err != nil {
			return err
		}

		if err := tx.Model(&book).Updates(newBook).Error; err != nil {
			return err
		}

		if newBook.Authors != nil {
			err := tx.Model(&book).Association("Authors").Replace(newBook.Authors)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func (br BookRepository) Delete(id uint) error {
	return br.db.Select("Authors").Delete(&models.Book{}, id).Error
}
