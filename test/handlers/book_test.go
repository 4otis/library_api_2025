package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/4otis/library_api_2025/internal/handlers"
	"github.com/4otis/library_api_2025/internal/migrations"
	"github.com/4otis/library_api_2025/internal/models"
	"github.com/4otis/library_api_2025/internal/repository"
	testutils "github.com/4otis/library_api_2025/test"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func setupBookHandler(t *testing.T) (*echo.Echo, *gorm.DB) {
	e := echo.New()
	db := testutils.SetupTestDB(t)
	err := migrations.RunInitMigrations(db)
	if err != nil {
		t.Fatal("Error. Failed to run InitMigrations.")
	}

	handlers.SetupRoutes(e, db)

	return e, db
}

func TestCreateBookHandler(t *testing.T) {
	e, db := setupBookHandler(t)
	defer testutils.FreeTestDB(t, db)

	t.Run("Create Book - Success", func(t *testing.T) {
		book := &models.Book{
			Title: "book1",
			Pages: 100,
		}
		body, _ := json.Marshal(book)

		req := httptest.NewRequest(http.MethodPost, "/books", bytes.NewReader((body)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)

		var resp models.Book
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		assert.NotZero(t, resp.ID)
		assert.Equal(t, book.Title, resp.Title)
		assert.Equal(t, book.Pages, resp.Pages)
	})

	t.Run("Create Book - Empty Title", func(t *testing.T) {
		book := &models.Book{
			Pages: 100,
		}
		body, _ := json.Marshal(book)

		req := httptest.NewRequest(http.MethodPost, "/books", bytes.NewReader((body)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)
	})

	t.Run("Create Book - Invalid ID (already existing)", func(t *testing.T) {
		book := &models.Book{
			Title: "book3",
			Pages: 300,
			Model: gorm.Model{
				ID: 1,
			},
		}
		body, _ := json.Marshal(book)

		req := httptest.NewRequest(http.MethodPost, "/books", bytes.NewReader((body)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("Create Book - With Authors", func(t *testing.T) {
		author1 := &models.Author{
			Name: "author1",
		}
		author2 := &models.Author{
			Name: "author2",
		}

		book := &models.Book{
			Title:   "book_with_authors",
			Pages:   150,
			Authors: []*models.Author{author1, author2},
		}
		body, _ := json.Marshal(book)

		req := httptest.NewRequest(http.MethodPost, "/books", bytes.NewReader((body)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)

		var resp models.Book
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		assert.NotZero(t, resp.ID)
		assert.Equal(t, book.Title, resp.Title)

		var cnt int64
		db.Table("books_authors").Where("book_id = ?", resp.ID).Count(&cnt)
		assert.Equal(t, int64(2), cnt)
	})
}

func TestGetBookHandler(t *testing.T) {
	e, db := setupBookHandler(t)
	defer testutils.FreeTestDB(t, db)

	t.Run("Get Book - Success", func(t *testing.T) {
		book := &models.Book{
			Title: "book1",
			Pages: 100,
		}
		body, _ := json.Marshal(book)

		req := httptest.NewRequest(http.MethodPost, "/books", bytes.NewReader((body)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		req = httptest.NewRequest(http.MethodGet, "/books/1", nil)
		rec = httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		var resp models.Book
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))

		assert.Equal(t, book.Title, resp.Title)
		assert.Equal(t, book.Pages, resp.Pages)
	})

	t.Run("Get Book - Invalid ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/books/999", nil)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}

func TestListBooksHandler(t *testing.T) {
	e, db := setupBookHandler(t)
	defer testutils.FreeTestDB(t, db)

	a1 := &models.Author{
		Name: "a1",
	}
	a2 := &models.Author{
		Name: "a2",
	}
	a3 := &models.Author{
		Name: "a3",
	}

	bookRepo := repository.NewBookRepository(db)
	books := []*models.Book{
		{
			Title:   "b1",
			Pages:   100,
			Authors: []*models.Author{a1},
		},
		{
			Title:   "b2",
			Pages:   200,
			Authors: []*models.Author{a1, a2},
		},
		{
			Title:   "b3",
			Pages:   300,
			Authors: []*models.Author{a2, a3},
		},
	}

	for _, book := range books {
		require.NoError(t, bookRepo.Create(book))
	}

	t.Run("List Books - Success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/books", nil)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response []models.Book
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response))

		assert.Len(t, response, 3, "excepted 3 books")
		for i := range 3 {
			assert.Equal(t, books[i].Title, response[i].Title)
			assert.Equal(t, books[i].Pages, response[i].Pages)
		}

		var cnt int64
		db.Table("books_authors").Count(&cnt)
		assert.Equal(t, int64(5), cnt, "excepted 5 notes in many2many table")
	})

	t.Run("List Books - Empty list", func(t *testing.T) {
		db.Exec("delete from books_authors")
		db.Exec("delete from authors")
		db.Exec("delete from books")

		req := httptest.NewRequest(http.MethodGet, "/books", nil)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response []models.Book
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response))

		assert.Empty(t, response, "response should be empty")
	})
}

func TestUpdateBookHandler(t *testing.T) {
	e, db := setupBookHandler(t)
	defer testutils.FreeTestDB(t, db)

	a1 := &models.Author{
		Name: "a1",
	}
	a2 := &models.Author{
		Name: "a2",
	}

	bookRepo := repository.NewBookRepository(db)
	books := []*models.Book{
		{
			Title:   "b1",
			Pages:   100,
			Authors: []*models.Author{a1},
		},
		{
			Title:   "b2",
			Pages:   200,
			Authors: []*models.Author{a2},
		},
	}

	require.NoError(t, bookRepo.Create(books[0]))

	t.Run("Update Book - Success", func(t *testing.T) {
		newBook := books[1]
		id := 1

		body, _ := json.Marshal(newBook)

		req := httptest.NewRequest(http.MethodPut, "/books/"+strconv.Itoa(id), bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNoContent, rec.Code)

		updatedBook, err := bookRepo.Read(uint(id))
		require.NoError(t, err, "failed to read updated book")

		assert.Equal(t, newBook.Title, updatedBook.Title)
		assert.Equal(t, newBook.Pages, updatedBook.Pages)
	})

	t.Run("Update Book - Partial update", func(t *testing.T) {
		id := 1
		newBook := books[0]
		newBook.Title = books[1].Title
		newBook.Pages = books[1].Pages

		body, _ := json.Marshal(newBook)

		req := httptest.NewRequest(http.MethodPut, "/books/"+strconv.Itoa(id), bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNoContent, rec.Code)

		updatedBook, err := bookRepo.Read(uint(id))
		require.NoError(t, err, "failed to read updated book")

		assert.Equal(t, newBook.Title, updatedBook.Title)
		assert.Equal(t, newBook.Pages, updatedBook.Pages)
	})

	t.Run("Update Book - Invalid ID (not found)", func(t *testing.T) {
		id := 999
		body, _ := json.Marshal(books[0])

		req := httptest.NewRequest(http.MethodPut, "/books/"+strconv.Itoa(id), bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}

func TestDeleteBookHandler(t *testing.T) {
	e, db := setupBookHandler(t)
	defer testutils.FreeTestDB(t, db)

	a1 := &models.Author{
		Name: "a1",
	}
	a2 := &models.Author{
		Name: "a2",
	}

	bookRepo := repository.NewBookRepository(db)
	books := []*models.Book{
		{
			Title:   "b1",
			Pages:   100,
			Authors: []*models.Author{a1, a2},
		},
		{
			Title:   "b2",
			Pages:   200,
			Authors: []*models.Author{a2},
		},
	}

	for _, book := range books {
		require.NoError(t, bookRepo.Create(book))
	}

	t.Run("Delete Book - Success", func(t *testing.T) {
		id := 1

		req := httptest.NewRequest(http.MethodDelete, "/books/"+strconv.Itoa(id), nil)
		rec := httptest.NewRecorder()

		var cnt int64
		db.Table("books").Count(&cnt)
		assert.Equal(t, int64(2), cnt)

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNoContent, rec.Code)

		_, err := bookRepo.Read(uint(id))
		require.Error(t, err)
	})
}
