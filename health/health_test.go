//go:build unit

package health

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestHealthHandler(t *testing.T) {

	req, err := http.NewRequest(http.MethodGet, "/health", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	resperr := HealthHandler(c)

	if assert.NoError(t, resperr) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "healthy", rec.Body.String())
	}
}
