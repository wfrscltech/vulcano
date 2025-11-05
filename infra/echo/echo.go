package echo

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	echom "github.com/labstack/echo/v4/middleware"
	"github.com/wfrscltech/vulcano/infra/echo/middleware"
)

// NewEchoInstance Crea e inicializa una nueva instancia de Echo
func NewEchoInstance(logger *slog.Logger, version, buildTime, commitHash string) *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Logger.SetOutput(os.Stdout)

	e.Use(echom.RequestLoggerWithConfig(echom.RequestLoggerConfig{
		LogStatus:   true,
		LogURI:      true,
		LogError:    true,
		HandleError: true, // forwards error to the global error handler, so it can decide appropriate status code
		LogValuesFunc: func(c echo.Context, v echom.RequestLoggerValues) error {
			if v.Error == nil {
				logger.LogAttrs(context.Background(), slog.LevelInfo, "REQUEST",
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
				)
			} else {
				logger.LogAttrs(context.Background(), slog.LevelError, "REQUEST_ERROR",
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
					slog.String("err", v.Error.Error()),
				)
			}
			return nil
		},
	}))
	e.Use(middleware.ProblemMiddleware)
	e.Use(echom.CORSWithConfig(echom.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodOptions, http.MethodPut},
		AllowHeaders: []string{"*"},
	}))

	slog.SetDefault(logger)

	hh := NewHealthHandler(version, buildTime, commitHash)
	e.GET("/health", hh.Healthcheck)

	return e
}
