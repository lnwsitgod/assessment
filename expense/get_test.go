//go:build unit

package expense

import (
	"database/sql"
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

func TestGetExpenseHandler(t *testing.T) {
	t.Run("Test case for successful get expense by ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/expenses", strings.NewReader(""))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)
		c.SetPath("/:id")
		c.SetParamNames("id")
		c.SetParamValues("1")

		mockSql := "SELECT id, title, amount, note, tags FROM expenses WHERE id = $1"
		mockRows := sqlmock.NewRows([]string{"id", "title", "amount", "note", "tags"}).AddRow(1, "title", 100.0, "note", pq.Array([]string{"tag1", "tag2"}))
		mockDB, mock, err := sqlmock.New()

		db = mockDB
		mock.ExpectPrepare(regexp.QuoteMeta(mockSql)).ExpectQuery().WithArgs("1").WillReturnRows(mockRows)
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer mockDB.Close()

		err = GetExpenseHandler(c)

		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, `{"id":1,"title":"title","amount":100,"note":"note","tags":["tag1","tag2"]}`, strings.TrimSpace(rec.Body.String()))
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("Test case for unable to prepare query expense statement", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/expenses", strings.NewReader(""))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)
		c.SetPath("/:id")
		c.SetParamNames("id")
		c.SetParamValues("1")

		mockSql := "SELECT id, title, amount, note, tags FROM expenses WHERE id = $1"
		mockRows := sqlmock.NewRows([]string{"id", "title", "amount", "note", "tags"}).AddRow(1, "title", 100.0, "note", pq.Array([]string{"tag1", "tag2"}))
		mockDB, mock, err := sqlmock.New()

		db = mockDB
		mock.ExpectQuery(regexp.QuoteMeta(mockSql)).WithArgs(1).WillReturnRows(mockRows)
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer mockDB.Close()

		err = GetExpenseHandler(c)

		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusInternalServerError, rec.Code)
			assert.Equal(t, `{"message":"cannot prepare query expense statment"}`, strings.TrimSpace(rec.Body.String()))
		}
	})

	t.Run("Test case for getting expense by ID not found", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/expenses", strings.NewReader(""))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)
		c.SetPath("/:id")
		c.SetParamNames("id")
		c.SetParamValues("1")

		mockSql := "SELECT id, title, amount, note, tags FROM expenses WHERE id = $1"
		mockDB, mock, err := sqlmock.New()

		db = mockDB
		mock.ExpectPrepare(regexp.QuoteMeta(mockSql)).ExpectQuery().WithArgs("1").WillReturnError(sql.ErrNoRows)
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer mockDB.Close()

		err = GetExpenseHandler(c)

		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusNotFound, rec.Code)
			assert.Equal(t, `{"message":"expense not found"}`, strings.TrimSpace(rec.Body.String()))
		}
	})

	t.Run("Test case for unable to scan expense", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/expenses", strings.NewReader(""))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)
		c.SetPath("/:id")
		c.SetParamNames("id")
		c.SetParamValues("1")

		mockSql := "SELECT id, title, amount, note, tags FROM expenses WHERE id = $1"
		mockDB, mock, err := sqlmock.New()

		db = mockDB
		mock.ExpectPrepare(regexp.QuoteMeta(mockSql)).ExpectQuery().WithArgs(1).WillReturnError(sql.ErrNoRows)
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer mockDB.Close()

		err = GetExpenseHandler(c)

		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusInternalServerError, rec.Code)
			assert.Equal(t, `{"message":"unable to scan expense"}`, strings.TrimSpace(rec.Body.String()))
		}
	})
}

func TestGetExpensesHandler(t *testing.T) {
	t.Run("Test case for successful get all expenses", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/expenses", strings.NewReader(""))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)

		mockSql := "SELECT id, title, amount, note, tags FROM expenses"
		mockRows := sqlmock.NewRows([]string{"id", "title", "amount", "note", "tags"}).
			AddRow("1", "title1", 100.0, "note1", pq.Array([]string{"tag1", "tag2"})).
			AddRow("2", "title2", 200.0, "note2", pq.Array([]string{"tag11", "tag22"}))
		mockDB, mock, err := sqlmock.New()

		db = mockDB
		mock.ExpectQuery(regexp.QuoteMeta(mockSql)).WillReturnRows(mockRows)
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer mockDB.Close()

		err = GetExpensesHandler(c)

		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, `[{"id":1,"title":"title1","amount":100,"note":"note1","tags":["tag1","tag2"]},{"id":2,"title":"title2","amount":200,"note":"note2","tags":["tag11","tag22"]}]`, strings.TrimSpace(rec.Body.String()))
		}
	})

	t.Run("Test case for unable to prepare query all expense", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/expenses", strings.NewReader(""))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)

		mockSql := "SELECT id, title, amount, note, tags FROM expensesx"
		mockRows := sqlmock.NewRows([]string{"id", "title", "amount", "note", "tags"}).
			AddRow("1", "title1", 100.0, "note1", pq.Array([]string{"tag1", "tag2"})).
			AddRow("2", "title2", 200.0, "note2", pq.Array([]string{"tag11", "tag22"}))
		mockDB, mock, err := sqlmock.New()

		db = mockDB
		mock.ExpectQuery(regexp.QuoteMeta(mockSql)).WillReturnRows(mockRows)
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer mockDB.Close()

		err = GetExpensesHandler(c)

		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusInternalServerError, rec.Code)
			assert.Equal(t, `{"message":"cannot query expense"}`, strings.TrimSpace(rec.Body.String()))
		}
	})

	t.Run("Test case for unable to prepare scan all expense", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/expenses", strings.NewReader(""))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)

		mockSql := "SELECT id, title, amount, note, tags FROM expenses"
		mockRows := sqlmock.NewRows([]string{"id", "title"}).AddRow("invalid", "title invalid")
		mockDB, mock, err := sqlmock.New()

		db = mockDB
		mock.ExpectQuery(regexp.QuoteMeta(mockSql)).WillReturnRows(mockRows)
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer mockDB.Close()

		err = GetExpensesHandler(c)

		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusInternalServerError, rec.Code)
			assert.Equal(t, `{"message":"unable to scan expense"}`, strings.TrimSpace(rec.Body.String()))
		}
	})
}
