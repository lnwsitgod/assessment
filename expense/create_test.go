//go:build unit

package expense

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestCreateExpenseHandler(t *testing.T) {
	t.Run("Test case for successful creation of expense", func(t *testing.T) {
		body := `{"title":"title","amount":100,"note":"note","tags":["tag1","tag2"]}`
		req := httptest.NewRequest(http.MethodPost, "/expenses", strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)

		mockSql := "INSERT INTO expenses (title, amount, note, tags) values ($1, $2, $3, $4) RETURNING id"
		mockRows := sqlmock.NewRows([]string{"id"}).AddRow("1")
		mockDB, mock, err := sqlmock.New()

		db = mockDB
		mock.ExpectQuery(regexp.QuoteMeta(mockSql)).WithArgs("title", 100.0, "note", pq.Array([]string{"tag1", "tag2"})).WillReturnRows(mockRows)
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer mockDB.Close()

		err = CreateExpenseHandler(c)

		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusCreated, rec.Code)
			assert.Equal(t, `{"id":1,"title":"title","amount":100,"note":"note","tags":["tag1","tag2"]}`, strings.TrimSpace(rec.Body.String()))
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("Test case for failed creation of expense", func(t *testing.T) {
		body := `invalid request`
		req := httptest.NewRequest(http.MethodPost, "/expenses", strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)

		err := CreateExpenseHandler(c)

		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Equal(t, `{"message":"invalid request"}`, strings.TrimSpace(rec.Body.String()))
		}
	})

	t.Run("Test case for invalid request with empty title", func(t *testing.T) {
		body := `{"title":"","amount":100,"note":"note","tags":["tag1","tag2"]}`
		req := httptest.NewRequest(http.MethodPost, "/expenses", strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)

		err := CreateExpenseHandler(c)

		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Equal(t, `{"message":"title is required"}`, strings.TrimSpace(rec.Body.String()))
		}
	})

	t.Run("Test case for invalid request with negative amount", func(t *testing.T) {
		body := `{"title":"title","amount":0,"note":"note","tags":["tag1","tag2"]}`
		req := httptest.NewRequest(http.MethodPost, "/expenses", strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)

		err := CreateExpenseHandler(c)

		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Equal(t, `{"message":"amount is required and must be greater than 0"}`, strings.TrimSpace(rec.Body.String()))
		}
	})

	t.Run("Test case for invalid request with empty note", func(t *testing.T) {
		body := `{"title":"title","amount":10,"note":"","tags":["tag1","tag2"]}`
		req := httptest.NewRequest(http.MethodPost, "/expenses", strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)

		err := CreateExpenseHandler(c)

		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Equal(t, `{"message":"note is required"}`, strings.TrimSpace(rec.Body.String()))
		}
	})

	t.Run("Test case for invalid request with no tags", func(t *testing.T) {
		body := `{"title":"title","amount":10,"note":"note","tags":[]}`
		req := httptest.NewRequest(http.MethodPost, "/expenses", strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)

		err := CreateExpenseHandler(c)

		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Equal(t, `{"message":"at least one tag is required"}`, strings.TrimSpace(rec.Body.String()))
		}
	})

	t.Run("Test case for database error during creation of expense", func(t *testing.T) {
		body := `{"title":"title","amount":100,"note":"note","tags":["tag1","tag2"]}`
		req := httptest.NewRequest(http.MethodPost, "/expenses", strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)

		mockSql := "INSERT INTO expenses (title, amount, note, tags) values ($1, $2, $3, $4) RETURNING id"
		mockDB, mock, err := sqlmock.New()
		db = mockDB

		mock.ExpectQuery(regexp.QuoteMeta(mockSql)).WithArgs("title", 100.0, "note", pq.Array([]string{"tag1", "tag2"})).WillReturnError(errors.New("database error"))
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer mockDB.Close()

		err = CreateExpenseHandler(c)

		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusInternalServerError, rec.Code)
			assert.Equal(t, `{"message":"cannot insert data"}`, strings.TrimSpace(rec.Body.String()))
		}
	})

	t.Run("Test case for validating required fields during creation of expense", func(t *testing.T) {
		tests := []struct {
			name string
			e    Expense
			err  string
		}{
			{
				name: "Test case 1: empty title",
				e:    Expense{Title: "", Amount: 100, Note: "note", Tags: []string{"tag1", "tag2"}},
				err:  "title is required",
			},
			{
				name: "Test case 2: negative amount",
				e:    Expense{Title: "title", Amount: -100, Note: "note", Tags: []string{"tag1", "tag2"}},
				err:  "amount is required and must be greater than 0",
			},
			{
				name: "Test case 3: empty note",
				e:    Expense{Title: "title", Amount: 100, Note: "", Tags: []string{"tag1", "tag2"}},
				err:  "note is required",
			},
			{
				name: "Test case 4: no tags",
				e:    Expense{Title: "title", Amount: 100, Note: "note", Tags: []string{}},
				err:  "at least one tag is required",
			},
			{
				name: "Test case 5: valid expense",
				e:    Expense{Title: "title", Amount: 100, Note: "note", Tags: []string{"tag1", "tag2"}},
				err:  "",
			},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				err := test.e.Validate()
				if test.err == "" {
					assert.NoError(t, err)
				} else {
					assert.EqualError(t, err, test.err)
				}
			})
		}
	})
}
