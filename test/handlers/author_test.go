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

func setupAuthorHandler(t *testing.T) (*echo.Echo, *gorm.DB) {
	e := echo.New()
	db := testutils.SetupTestDB(t)
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

func TestCreateAuthorHandler(t *testing.T) {
	e, db := setupAuthorHandler(t)
	defer testutils.FreeTestDB(t, db)

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

func TestGetAuthorHandler(t *testing.T) {
	e, db := setupAuthorHandler(t)
	defer testutils.FreeTestDB(t, db)

	t.Run("Get Author - Success", func(t *testing.T) {
		author := &models.Author{
			Name: "author1",
		}
		body, _ := json.Marshal(author)

		req := httptest.NewRequest(http.MethodPost, "/authors", bytes.NewReader((body)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		req = httptest.NewRequest(http.MethodGet, "/authors/1", nil)
		rec = httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		var resp models.Author
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))

		assert.Equal(t, author.Name, resp.Name)
	})

	t.Run("Get Author - Invalid ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/authors/999", nil)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}

func TestListAuthorHandler(t *testing.T) {
	e, db := setupAuthorHandler(t)
	defer testutils.FreeTestDB(t, db)

	b1 := &models.Book{
		Title: "b1",
		Pages: 100,
	}
	b2 := &models.Book{
		Title: "b2",
		Pages: 200,
	}
	b3 := &models.Book{
		Title: "b3",
		Pages: 300,
	}

	authorRepo := repository.NewAuthorRepository(db)
	authors := []*models.Author{
		{
			Name:  "a1",
			Books: []*models.Book{b1},
		},
		{
			Name:  "a2",
			Books: []*models.Book{b1, b2},
		},
		{
			Name:  "a3",
			Books: []*models.Book{b2, b3},
		},
	}

	for _, author := range authors {
		require.NoError(t, authorRepo.Create(author))
	}

	t.Run("List Author - Success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/authors", nil)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response []models.Author
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response))

		assert.Len(t, response, 3, "excepted 3 authors")
		for i := range 3 {
			assert.Equal(t, authors[i].Name, response[i].Name)
			assert.Equal(t, authors[i].Books[0].Title, response[i].Books[0].Title)
		}

		var cnt int64
		db.Table("books_authors").Count(&cnt)
		assert.Equal(t, int64(5), cnt, "excepted 5 notes in many2many table")
	})

	t.Run("List Author - Empty list", func(t *testing.T) {
		db.Exec("delete from books_authors")
		db.Exec("delete from books")
		db.Exec("delete from authors")

		req := httptest.NewRequest(http.MethodGet, "/authors", nil)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response []models.Author
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response))

		assert.Empty(t, response, "response should be empty")
	})

}

func TestUpdateAuthorHandler(t *testing.T) {
	e, db := setupAuthorHandler(t)
	defer testutils.FreeTestDB(t, db)

	b1 := &models.Book{
		Title: "b1",
		Pages: 100,
	}
	b2 := &models.Book{
		Title: "b2",
		Pages: 200,
	}

	authorRepo := repository.NewAuthorRepository(db)
	authors := []*models.Author{
		{
			Name:  "a1",
			Books: []*models.Book{b1},
		},
		{
			Name:  "a2",
			Books: []*models.Book{b2},
		},
	}

	require.NoError(t, authorRepo.Create(authors[0]))

	t.Run("Update Author - Success", func(t *testing.T) {
		newAuthor := authors[1]
		id := 1

		body, _ := json.Marshal(newAuthor)

		req := httptest.NewRequest(http.MethodPut, "/authors/"+strconv.Itoa(id), bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNoContent, rec.Code)

		updatedAuthor, err := authorRepo.Read(uint(id))
		require.NoError(t, err, "failed to read updated author")

		assert.Equal(t, newAuthor.Name, updatedAuthor.Name)
		assert.Equal(t, newAuthor.Books[0].Title, updatedAuthor.Books[0].Title)

	})

	t.Run("Update Author - Partial update", func(t *testing.T) {
		id := 1
		newAuthor := authors[0]
		newAuthor.Name = authors[1].Name

		body, _ := json.Marshal(newAuthor)

		req := httptest.NewRequest(http.MethodPut, "/authors/"+strconv.Itoa(id), bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNoContent, rec.Code)

		updatedAuthor, err := authorRepo.Read(uint(id))
		require.NoError(t, err, "failed to read updated author")

		assert.Equal(t, newAuthor.Name, updatedAuthor.Name)
		assert.Equal(t, newAuthor.Books[0].Title, updatedAuthor.Books[0].Title)

	})

	t.Run("Update Author - Invalid ID (not found)", func(t *testing.T) {
		id := 999
		body, _ := json.Marshal(authors[0])

		req := httptest.NewRequest(http.MethodPut, "/authors/"+strconv.Itoa(id), bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)

	})

}

func TestDeleteAuthorHandler(t *testing.T) {
	e, db := setupAuthorHandler(t)
	defer testutils.FreeTestDB(t, db)

	b1 := &models.Book{
		Title: "b1",
		Pages: 100,
	}
	b2 := &models.Book{
		Title: "b2",
		Pages: 200,
	}

	authorRepo := repository.NewAuthorRepository(db)
	authors := []*models.Author{
		{
			Name:  "a1",
			Books: []*models.Book{b1, b2},
		},
		{
			Name:  "a2",
			Books: []*models.Book{b2},
		},
	}

	for _, author := range authors {
		require.NoError(t, authorRepo.Create(author))
	}

	t.Run("Delete Author - Success", func(t *testing.T) {
		id := 1

		req := httptest.NewRequest(http.MethodDelete, "/authors/"+strconv.Itoa(id), nil)
		rec := httptest.NewRecorder()

		var cnt int64
		db.Table("authors").Count(&cnt)
		assert.Equal(t, int64(2), cnt)

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNoContent, rec.Code)

		_, err := authorRepo.Read(uint(id))
		require.Error(t, err)
	})
}
