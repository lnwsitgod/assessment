package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/lnwsitgod/assessment/health"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/health", health.GetHealthHandler)

	g := e.Group("/expenses")
	g.Use(authMiddlewareGuard)
	g.GET("", func(c echo.Context) error {
		return c.String(http.StatusOK, "You're authenticated!")
	})

	startServerGracefullyShutdown(e)
}

func serverPort() string {
	port := "2565"
	if p, ok := os.LookupEnv("PORT"); ok {
		port = p
	}
	return fmt.Sprintf(":%s", port)
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
		if err := e.Start(serverPort()); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt)
	<-shutdown

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
