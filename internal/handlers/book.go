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

func (bh BookHandler) ListBooks(c echo.Context) error {
	books, err := bh.repository.ReadAll()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, books)
}

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
