package handlers

import (
	// _ "github.com/4otis/library_api_2025/docs"
	"github.com/4otis/library_api_2025/internal/repository"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	"gorm.io/gorm"
)

func SetupRoutes(e *echo.Echo, db *gorm.DB) {
	bookRepo := repository.NewBookRepository(db)
	authorRepo := repository.NewAuthorRepository(db)

	bookHandler := NewBookHandler(bookRepo)
	authorHandler := NewAuthorHandler(authorRepo)

	e.GET("/books", bookHandler.ListBooks)
	e.GET("/books/:id", bookHandler.GetBook)
	e.POST("/books", bookHandler.CreateBook)
	e.PUT("/books/:id", bookHandler.UpdateBook)
	e.DELETE("/books/:id", bookHandler.DeleteBook)

	e.GET("/authors", authorHandler.ListAuthors)
	e.GET("/authors/:id", authorHandler.GetAuthor)
	e.POST("/authors", authorHandler.CreateAuthor)
	e.PUT("/authors/:id", authorHandler.UpdateAuthor)
	e.DELETE("/authors/:id", authorHandler.DeleteAuthor)

	e.GET("/swagger/*", echoSwagger.WrapHandler)
}
