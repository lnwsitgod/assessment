package health

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetHealthHandler(c echo.Context) error {
	return c.String(http.StatusOK, "healthy")
}
