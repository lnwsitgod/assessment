//go:build integration

package expense

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestIntegrationCreateExpenseHandler(t *testing.T) {
	InitDB()
	defer CloseDB()

	teardown := startIntegrationTestServer(t)
	defer teardown()

	var ep Expense
	body := bytes.NewBufferString(`{"title":"TestIntegrationCreateExpenseHandler","amount":100,"note":"TestIntegrationCreateExpenseHandler note","tags":["integration","test", "create"]}`)

	res := request(http.MethodPost, uri("expenses"), body)
	err := res.Decode(&ep)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, res.StatusCode)
	assert.NotEqual(t, 0, ep.ID)
	assert.Equal(t, "TestIntegrationCreateExpenseHandler", ep.Title)
	assert.Equal(t, 100.0, ep.Amount)
	assert.Equal(t, "TestIntegrationCreateExpenseHandler note", ep.Note)
	assert.Equal(t, "integration", ep.Tags[0])
	assert.Equal(t, "test", ep.Tags[1])
	assert.Equal(t, []string([]string{"integration", "test", "create"}), ep.Tags)
}

func TestIntegrationGetExpenseHandler(t *testing.T) {
	InitDB()
	defer CloseDB()

	teardown := startIntegrationTestServer(t)
	defer teardown()

	e := seedExpense(t)
	var ep Expense

	res := request(http.MethodGet, uri("expenses", strconv.Itoa(e.ID)), nil)
	err := res.Decode(&ep)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.NotEqual(t, 0, ep.ID)
	assert.Equal(t, "integration test title", ep.Title)
	assert.Equal(t, 100.0, ep.Amount)
	assert.Equal(t, "integration test note", ep.Note)
	assert.Equal(t, "integration", ep.Tags[0])
	assert.Equal(t, "test", ep.Tags[1])
	assert.Equal(t, []string([]string{"integration", "test"}), ep.Tags)
}

func TestIntegrationUpdateExpenseHandler(t *testing.T) {
	InitDB()
	defer CloseDB()

	teardown := startIntegrationTestServer(t)
	defer teardown()

	e := seedExpense(t)
	var ep Expense
	body := bytes.NewBufferString(`{"title":"TestIntegrationUpdateExpenseHandler","amount":100,"note":"TestIntegrationUpdateExpenseHandler note","tags":["integration","test", "update"]}`)

	res := request(http.MethodPut, uri("expenses", strconv.Itoa(e.ID)), body)
	err := res.Decode(&ep)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.NotEqual(t, 0, ep.ID)
	assert.Equal(t, "TestIntegrationUpdateExpenseHandler", ep.Title)
	assert.Equal(t, 100.0, ep.Amount)
	assert.Equal(t, "TestIntegrationUpdateExpenseHandler note", ep.Note)
	assert.Equal(t, "integration", ep.Tags[0])
	assert.Equal(t, "test", ep.Tags[1])
	assert.Equal(t, []string([]string{"integration", "test", "update"}), ep.Tags)
}

func TestIntegrationGetExpensesHandler(t *testing.T) {
	InitDB()
	defer CloseDB()

	teardown := startIntegrationTestServer(t)
	defer teardown()

	seedExpense(t)
	var eps []Expense

	res := request(http.MethodGet, uri("expenses"), nil)
	err := res.Decode(&eps)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Greater(t, len(eps), 0)
}

func startIntegrationTestServer(t *testing.T) func() {
	e := echo.New()

	go func() {
		e.Use(authMiddlewareGuardIntegrationTest)

		e.POST("/expenses", CreateExpenseHandler)
		e.GET("/expenses/:id", GetExpenseHandler)
		e.PUT("/expenses/:id", UpdateExpenseHandler)
		e.GET("expenses", GetExpensesHandler)
		e.Start(fmt.Sprintf(":%s", os.Getenv("PORT")))
	}()
	for {
		conn, err := net.DialTimeout("tcp", fmt.Sprintf("localhost:%s", os.Getenv("PORT")), 30*time.Second)
		if err != nil {
			log.Println(err)
		}
		if conn != nil {
			conn.Close()
			break
		}
	}

	return func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		err := e.Shutdown(ctx)
		assert.NoError(t, err)
	}
}

func request(method, url string, body io.Reader) *Response {
	req, _ := http.NewRequest(method, url, body)
	req.Header.Add("Authorization", os.Getenv("AUTH_TOKEN"))
	req.Header.Add("Content-Type", "application/json")
	client := http.Client{}
	res, err := client.Do(req)
	return &Response{res, err}
}

type Response struct {
	*http.Response
	err error
}

func authMiddlewareGuardIntegrationTest(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("Authorization")
		if token != os.Getenv("AUTH_TOKEN") {
			return c.String(http.StatusUnauthorized, "Unauthorized")
		}
		return next(c)
	}
}

func seedExpense(t *testing.T) Expense {
	var c Expense
	body := bytes.NewBufferString(`{"title":"integration test title","amount":100,"note":"integration test note","tags":["integration","test"]}`)
	err := request(http.MethodPost, uri("expenses"), body).Decode(&c)
	if err != nil {
		t.Fatal("can't create expense:", err)
	}
	return c
}

func uri(paths ...string) string {
	host := "http://localhost:2565"
	if paths == nil {
		return host
	}

	url := append([]string{host}, paths...)
	return strings.Join(url, "/")
}

func (r *Response) Decode(v interface{}) error {
	if r.err != nil {
		return r.err
	}

	return json.NewDecoder(r.Body).Decode(v)
}
