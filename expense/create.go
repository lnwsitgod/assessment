package expense

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

func CreateExpenseHandler(c echo.Context) error {
	e := Expense{}
	err := c.Bind(&e)
	if err != nil {
		c.Logger().Error("invalid request binding to struct exepnse error: ", err)
		return c.JSON(http.StatusBadRequest, Err{Message: errors.New("invalid request").Error()})
	}

	if err := e.Validate(); err != nil {
		c.Logger().Error("invalid request error: ", err)
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	row := db.QueryRow("INSERT INTO expenses (title, amount, note, tags) values ($1, $2, $3, $4) RETURNING id", e.Title, e.Amount, e.Note, pq.Array(e.Tags))
	err = row.Scan(&e.ID)
	if err != nil {
		c.Logger().Error("insert data error: ", err)
		return c.JSON(http.StatusInternalServerError, Err{Message: errors.New("cannot insert data").Error()})
	}

	return c.JSON(http.StatusCreated, e)
}
