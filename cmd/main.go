package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/labstack/echo/v4"
	stdmw "github.com/labstack/echo/v4/middleware"
)

func main() {
	if err := run(); err != nil {
		log.Println("Failed to run the server: ", err)
		os.Exit(1)
	}
}

func run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	e := echo.New()
	e.HideBanner = true
	e.HTTPErrorHandler = httpErr
	e.Use(stdmws()...)
	e.GET("/_healthz", func(c echo.Context) error {
		return c.JSON(http.StatusOK, echo.Map{
			"code":    http.StatusOK,
			"status":  "OK",
			"message": "Available!",
		})
	})

	errChan := make(chan error, 1)
	go func() {
		errChan <- e.Start(fmt.Sprintf(":%s", getEnv("PORT", "8280")))
	}()

	ctx, cancel = signal.NotifyContext(ctx, os.Interrupt, os.Kill, syscall.SIGTERM)
	defer cancel()

	select {
	case <-ctx.Done():
		log.Println("Shutting down the server...")
		if err := e.Shutdown(ctx); err != nil {
			return fmt.Errorf("failed to shutdown the server: %w", err)
		}
		log.Println("Server shut down")

	case err := <-errChan:
		if err != nil && err != http.ErrServerClosed {
			return err
		}
	}
	return nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func stdmws() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{
		stdmw.Recover(),
		stdmw.RemoveTrailingSlash(),
		stdmw.CORS(),
		stdmw.Logger(),
		stdmw.RateLimiter(stdmw.NewRateLimiterMemoryStore(10)),
		stdmw.Secure(),
	}
}

func httpErr(err error, c echo.Context) {
	c.Logger().Error("HTTP error", "err", err)

	if he, ok := err.(*echo.HTTPError); ok {
		switch he.Code {
		case http.StatusNotFound:
			c.JSON(http.StatusNotFound, echo.Map{
				"message": "Not found",
				"code":    http.StatusNotFound,
				"status":  "NOT_FOUND",
			})
			return

		case http.StatusTooManyRequests:
			c.JSON(http.StatusTooManyRequests, echo.Map{
				"message": "Too many requests. Please try again later.",
				"code":    http.StatusTooManyRequests,
				"status":  "RESOURCE_EXHAUSTED",
			})
			return

		default:
			c.JSON(he.Code, echo.Map{
				"message": "Unknown error",
				"code":    he.Code,
				"status":  "UNKNOWN_ERROR",
			})
		}
	}

	c.JSON(http.StatusInternalServerError, echo.Map{
		"message": "An internal error occurred",
		"code":    http.StatusInternalServerError,
		"status":  "INTERNAL_ERROR",
	})
}
