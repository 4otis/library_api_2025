package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/4otis/library_api_2025/internal/models"
	"github.com/4otis/library_api_2025/internal/repository"
	"github.com/labstack/echo/v4"
)

type AuthorHandler struct {
	repository *repository.AuthorRepository
}

func NewAuthorHandler(r *repository.AuthorRepository) *AuthorHandler {
	return &AuthorHandler{repository: r}
}

func (ah AuthorHandler) ListAuthors(c echo.Context) error {
	authors, err := ah.repository.ReadAll()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, authors)
}

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
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}

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
