package expense

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

func UpdateExpenseHandler(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Logger().Error("cast id error: ", err)
		return c.JSON(http.StatusBadRequest, Err{Message: errors.New("invalid request").Error()})
	}

	e := Expense{}
	err = c.Bind(&e)
	if err != nil {
		c.Logger().Error("invalid request binding to struct exepnse error: ", err)
		return c.JSON(http.StatusBadRequest, Err{Message: errors.New("invalid request").Error()})
	}

	stmt, err := db.Prepare("UPDATE expenses SET title = $2, amount = $3, note = $4, tags = $5 WHERE id = $1")
	if err != nil {
		c.Logger().Error("prepare statment error: ", err)
		return c.JSON(http.StatusInternalServerError, Err{Message: errors.New("cannot prepare query expense statment").Error()})
	}

	if _, err = stmt.Exec(id, e.Title, e.Amount, e.Note, pq.Array(e.Tags)); err != nil {
		c.Logger().Error("update data error: ", err)
		return c.JSON(http.StatusInternalServerError, Err{Message: errors.New("cannot update data").Error()})
	}
	e.ID = id
	return c.JSON(http.StatusOK, e)
}
