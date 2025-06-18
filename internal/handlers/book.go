package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/4otis/library_api_2025/internal/models"
	"github.com/4otis/library_api_2025/internal/repository"
	"github.com/labstack/echo/v4"
)

type BookHandler struct {
	repository *repository.BookRepository
}

func NewBookHandler(r *repository.BookRepository) *BookHandler {
	return &BookHandler{repository: r}
}

// ListBooks godoc
// @Summary Get all books
// @Description Get details of all books
// @Tags books
// @Accept  json
// @Produce  json
// @Success 200 {array} models.Book
// @Router /books [get]
func (bh BookHandler) ListBooks(c echo.Context) error {
	books, err := bh.repository.ReadAll()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, books)
}

// GetBook godoc
// @Summary Get book by ID
// @Description Get detailed information about a specific book
// @Tags books
// @Accept json
// @Produce json
// @Param id path int true "Book ID"
// @Success 200 {object} models.Book
// @Failure 400 {object} map[string]string "Invalid ID format"
// @Failure 404 {object} map[string]string "Book not found"
// @Router /books/{id} [get]
func (bh BookHandler) GetBook(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Error. Invalid ID format.")
	}

	book, err := bh.repository.Read(uint(id))
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("Error. Book not found (by id: %d).", id))
	}

	return c.JSON(http.StatusOK, book)
}

// CreateBook godoc
// @Summary Create a new book
// @Description Add a new book to the library
// @Tags books
// @Accept json
// @Produce json
// @Param book body models.Book true "Book data"
// @Success 201 {object} models.Book
// @Failure 400 {object} map[string]string "Invalid request body"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /books [post]
func (bh BookHandler) CreateBook(c echo.Context) error {
	var book models.Book
	err := c.Bind(&book)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Error. Invalid request body.")
	}

	err = bh.repository.Create(&book)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, book)
}

// UpdateBook godoc
// @Summary Update book information
// @Description Update existing book's data
// @Tags books
// @Accept json
// @Produce json
// @Param id path int true "Book ID"
// @Param book body models.Book true "Updated book data"
// @Success 204 "No content"
// @Failure 400 {object} map[string]string "Invalid ID format or request body"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /books/{id} [put]
func (bh BookHandler) UpdateBook(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Error. Invalid ID format.")
	}

	var book models.Book
	err = c.Bind(&book)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Error. Invalid request body.")
	}

	err = bh.repository.Update(uint(id), &book)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}

// DeleteBook godoc
// @Summary Delete a book
// @Description Remove book from the library
// @Tags books
// @Accept json
// @Produce json
// @Param id path int true "Book ID"
// @Success 204 "No content"
// @Failure 400 {object} map[string]string "Invalid ID format"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /books/{id} [delete]
func (bh BookHandler) DeleteBook(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Error. Invalid ID format.")
	}

	err = bh.repository.Delete(uint(id))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}
