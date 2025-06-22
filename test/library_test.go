package library_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/4otis/library_api_2025/internal/handlers"
	"github.com/4otis/library_api_2025/internal/migrations"
	"github.com/4otis/library_api_2025/internal/models"
	"github.com/4otis/library_api_2025/internal/repository"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func SetupTestDB(t *testing.T) *gorm.DB {
	dsn := "host=localhost user=postgres password=password dbname=library port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("Error. Failed to connect to test DB: %v", err)
	}

	schema := "test_" + t.Name()
	db.Exec("drop schema if exists " + schema + " cascade")
	db.Exec("create schema " + schema)
	db.Exec("set search_path to " + schema)

	return db
}

func FreeTestDB(t *testing.T, db *gorm.DB) {
	schema := "test_" + t.Name()
	db.Exec("drop schema if exists " + schema + " cascade")
}

func setupAuthorHandler(t *testing.T) (*echo.Echo, *gorm.DB) {
	e := echo.New()
	db := SetupTestDB(t)
	err := migrations.RunInitMigrations(db)
	if err != nil {
		t.Fatal("Error. Failed to run InitMigrations.")
	}
	repo := repository.NewAuthorRepository(db)
	handler := handlers.NewAuthorHandler(repo)

	e.POST("/authors", handler.CreateAuthor)
	e.GET("/authors", handler.ListAuthors)
	e.GET("/authors/:id", handler.GetAuthor)
	e.PUT("/authors/:id", handler.UpdateAuthor)
	e.DELETE("/authors/:id", handler.DeleteAuthor)

	return e, db
}

func TestAuthorCreateHandler(t *testing.T) {
	e, db := setupAuthorHandler(t)
	defer FreeTestDB(t, db)

	t.Run("Create Author - Success", func(t *testing.T) {
		author := &models.Author{
			Name: "author1",
		}
		body, _ := json.Marshal(author)

		req := httptest.NewRequest(http.MethodPost, "/authors", bytes.NewReader((body)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)

		var resp models.Author
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		assert.NotZero(t, resp.ID)
		assert.Equal(t, author.Name, resp.Name)
	})

	t.Run("Create Author - Empty Name", func(t *testing.T) {
		author := &models.Author{
			Name: "",
		}
		body, _ := json.Marshal(author)

		req := httptest.NewRequest(http.MethodPost, "/authors", bytes.NewReader((body)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)

		var resp models.Author
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		assert.Equal(t, author.Name, resp.Name)
		assert.NotZero(t, resp.ID)
	})

	t.Run("Create Author - Invalid ID (already existing)", func(t *testing.T) {
		author := &models.Author{
			Name: "author3",
			Model: gorm.Model{
				ID: 1,
			},
		}
		body, _ := json.Marshal(author)

		req := httptest.NewRequest(http.MethodPost, "/authors", bytes.NewReader((body)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("Create Author - Before book was added", func(t *testing.T) {
		book1 := &models.Book{
			Title: "book1",
			Model: gorm.Model{
				ID: 1,
			}}

		author := &models.Author{
			Name:  "author1",
			Books: []*models.Book{book1},
		}
		body, _ := json.Marshal(author)

		req := httptest.NewRequest(http.MethodPost, "/authors", bytes.NewReader((body)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)

		var resp models.Author
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		assert.NotZero(t, resp.ID)
		assert.Equal(t, author.Name, resp.Name)

		var cnt int64
		db.Table("books_authors").Where("author_id = ?", resp.ID).Count(&cnt)
		assert.Equal(t, int64(1), cnt)
	})

}
