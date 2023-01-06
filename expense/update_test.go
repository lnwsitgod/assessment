//go:build unit

package expense

import (
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

func TestUpdateExpenseHandler(t *testing.T) {
	t.Run("Test case for successful update expense by ID", func(t *testing.T) {
		body := `{"title":"update title","amount":99.9,"note":"note update","tags":["update1", "update2"]}`
		req := httptest.NewRequest(http.MethodPut, "/expenses", strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)
		c.SetPath("/:id")
		c.SetParamNames("id")
		c.SetParamValues("1")

		mockSql := "UPDATE expenses SET title = $2, amount = $3, note = $4, tags = $5 WHERE id = $1"
		mockDB, mock, err := sqlmock.New()

		db = mockDB
		mock.ExpectPrepare(regexp.QuoteMeta(mockSql)).ExpectExec().WithArgs(1, "update title", 99.9, "note update", pq.Array([]string{"update1", "update2"})).WillReturnResult(sqlmock.NewResult(0, 0))
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer mockDB.Close()

		err = UpdateExpenseHandler(c)

		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, `{"id":1,"title":"update title","amount":99.9,"note":"note update","tags":["update1","update2"]}`, strings.TrimSpace(rec.Body.String()))
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("Test case for failed update expense when casting id", func(t *testing.T) {
		body := `{"title":"update title","amount":99.9,"note":"note update","tags":["update1", "update2"]}`
		req := httptest.NewRequest(http.MethodPut, "/expenses", strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)
		c.SetPath("/:id")
		c.SetParamNames("id")
		c.SetParamValues("d")

		err := UpdateExpenseHandler(c)

		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Equal(t, `{"message":"invalid request"}`, strings.TrimSpace(rec.Body.String()))
		}
	})

	t.Run("Test case for failed update when bind value to struct expense", func(t *testing.T) {
		body := `invalid request`
		req := httptest.NewRequest(http.MethodPost, "/expenses", strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)
		c.SetPath("/:id")
		c.SetParamNames("id")
		c.SetParamValues("1")

		err := UpdateExpenseHandler(c)

		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Equal(t, `{"message":"invalid request"}`, strings.TrimSpace(rec.Body.String()))
		}
	})

	t.Run("Test case for unable to prepare update expense statement", func(t *testing.T) {
		body := `{"title":"update title","amount":99.9,"note":"note update","tags":["update1", "update2"]}`
		req := httptest.NewRequest(http.MethodPut, "/expenses", strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)
		c.SetPath("/:id")
		c.SetParamNames("id")
		c.SetParamValues("1")

		mockSql := "UPDATE expenses SET title = $2, amount = $3, note = $4, tags = $5 WHERE id = $1"
		mockDB, mock, err := sqlmock.New()

		db = mockDB
		mock.ExpectQuery(regexp.QuoteMeta(mockSql)).WithArgs(1)
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer mockDB.Close()

		err = UpdateExpenseHandler(c)

		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusInternalServerError, rec.Code)
			assert.Equal(t, `{"message":"cannot prepare query expense statment"}`, strings.TrimSpace(rec.Body.String()))
		}
	})

	t.Run("Test case for database error during update of expense", func(t *testing.T) {
		body := `{"title":"update title","amount":99.9,"note":"note update","tags":["update1", "update2"]}`
		req := httptest.NewRequest(http.MethodPut, "/expenses", strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)
		c.SetPath("/:id")
		c.SetParamNames("id")
		c.SetParamValues("1")

		mockSql := "UPDATE expenses SET title = $2, amount = $3, note = $4, tags = $5 WHERE id = $1"
		mockDB, mock, err := sqlmock.New()
		db = mockDB

		mock.ExpectPrepare(regexp.QuoteMeta(mockSql)).ExpectExec().WithArgs("1", "update title", 99.9, "note update", pq.Array([]string{"update1", "update2"})).WillReturnResult(sqlmock.NewResult(0, 0))
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer mockDB.Close()

		err = UpdateExpenseHandler(c)

		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusInternalServerError, rec.Code)
			assert.Equal(t, `{"message":"cannot update data"}`, strings.TrimSpace(rec.Body.String()))
		}
	})
}
