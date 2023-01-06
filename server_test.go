//go:build unit

package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestAuthGuardValidToken(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	assert.NoError(t, err)
	req.Header.Add("Authorization", "November 10, 2009")

	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	handler := func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	}
	chain := authMiddlewareGuard(handler)

	if assert.NoError(t, chain(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "OK", rec.Body.String())
	}
}

func TestAuthGuardInvalidToken(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	assert.NoError(t, err)
	req.Header.Add("Authorization", "November 10, 2009wrong_token")

	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	handler := func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	}
	chain := authMiddlewareGuard(handler)

	if assert.NoError(t, chain(c)) {
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		assert.Equal(t, "Unauthorized", rec.Body.String())
	}
}
