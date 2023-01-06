package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/lnwsitgod/assessment/expense"
	"github.com/lnwsitgod/assessment/health"
)

func main() {
	e := echo.New()
	e.Logger.SetLevel(log.INFO)
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/health", health.GetHealthHandler)

	g := e.Group("/expenses")
	g.Use(authMiddlewareGuard)
	g.POST("", expense.CreateExpenseHandler)
	g.GET("/:id", expense.GetExpenseHandler)
	g.PUT("/:id", expense.UpdateExpenseHandler)

	startServerGracefullyShutdown(e)
}

func authMiddlewareGuard(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("Authorization")
		if token != os.Getenv("AUTH_TOKEN") {
			return c.String(http.StatusUnauthorized, "Unauthorized")
		}
		return next(c)
	}
}

func startServerGracefullyShutdown(e *echo.Echo) {
	go func() {
		if err := e.Start(fmt.Sprintf(":%s", os.Getenv("PORT"))); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("starting server error:", err)
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	<-shutdown
	e.Logger.Info("shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	} else {
		e.Logger.Info("http server stopped")
	}
}
