package expense

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

func GetExpenseHandler(c echo.Context) error {
	id := c.Param("id")

	stmt, err := db.Prepare("SELECT id, title, amount, note, tags FROM expenses WHERE id = $1")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: errors.New("cannot prepare query expense statment").Error()})
	}

	row := stmt.QueryRow(id)
	e := Expense{}
	err = row.Scan(&e.ID, &e.Title, &e.Amount, &e.Note, pq.Array(&e.Tags))
	if err == sql.ErrNoRows {
		return c.JSON(http.StatusNotFound, Err{Message: errors.New("expense not found").Error()})
	} else if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: errors.New("unable to scan expense").Error()})
	}

	return c.JSON(http.StatusOK, e)
}
