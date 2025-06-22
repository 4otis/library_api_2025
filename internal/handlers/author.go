package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/4otis/library_api_2025/internal/models"
	"github.com/4otis/library_api_2025/internal/repository"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type AuthorHandler struct {
	repository *repository.AuthorRepository
}

func NewAuthorHandler(r *repository.AuthorRepository) *AuthorHandler {
	return &AuthorHandler{repository: r}
}

// ListAuthors godoc
// @Summary Get all authors
// @Description Get details of all authors
// @Tags authors
// @Accept  json
// @Produce  json
// @Success 200 {array} models.Author
// @Router /authors [get]
func (ah AuthorHandler) ListAuthors(c echo.Context) error {
	authors, err := ah.repository.ReadAll()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, authors)
}

// GetAuthor godoc
// @Summary Get author by ID
// @Description Get detailed information about a specific author
// @Tags authors
// @Accept json
// @Produce json
// @Param id path int true "Author ID"
// @Success 200 {object} models.Author
// @Failure 400 {object} map[string]string "Invalid ID format"
// @Failure 404 {object} map[string]string "Author not found"
// @Router /authors/{id} [get]
func (ah AuthorHandler) GetAuthor(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Error. Invalid ID format.")
	}

	author, err := ah.repository.Read(uint(id))
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("Error. Author not found (by id: %d).", id))
	}

	return c.JSON(http.StatusOK, author)
}

// CreateAuthor godoc
// @Summary Create a new author
// @Description Add a new author to the system
// @Tags authors
// @Accept json
// @Produce json
// @Param author body models.Author true "Author data"
// @Success 201 {object} models.Author
// @Failure 400 {object} map[string]string "Invalid request body"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /authors [post]
func (ah AuthorHandler) CreateAuthor(c echo.Context) error {
	var author models.Author
	err := c.Bind(&author)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Error. Invalid request body.")
	}

	err = ah.repository.Create(&author)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, author)
}

// UpdateAuthor godoc
// @Summary Update author information
// @Description Update existing author's data
// @Tags authors
// @Accept json
// @Produce json
// @Param id path int true "Author ID"
// @Param author body models.Author true "Updated author data"
// @Success 204 "No content"
// @Failure 400 {object} map[string]string "Invalid ID format or request body"
// @Failure 404 {object} map[string]string "Author not found by entered id"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /authors/{id} [put]
func (ah AuthorHandler) UpdateAuthor(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Error. Invalid ID format.")
	}

	var author models.Author
	err = c.Bind(&author)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Error. Invalid request body.")
	}

	err = ah.repository.Update(uint(id), &author)
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("Error. Author not found (by id: %d).", id))
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}
	return c.NoContent(http.StatusNoContent)
}

// DeleteAuthor godoc
// @Summary Delete an author
// @Description Remove author from the system
// @Tags authors
// @Accept json
// @Produce json
// @Param id path int true "Author ID"
// @Success 204 "No content"
// @Failure 400 {object} map[string]string "Invalid ID format"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /authors/{id} [delete]
func (ah AuthorHandler) DeleteAuthor(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Error. Invalid ID format.")
	}

	err = ah.repository.Delete(uint(id))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}
