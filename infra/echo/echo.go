package echo

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	echom "github.com/labstack/echo/v4/middleware"
	"github.com/wfrscltech/vulcano/infra/echo/apidocs"
	"github.com/wfrscltech/vulcano/infra/echo/middleware"
)

// NewEchoInstance Crea e inicializa una nueva instancia de Echo
func NewEchoInstance(logger *slog.Logger, version, buildTime, commitHash string) *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Logger.SetOutput(os.Stdout)

	e.Use(middleware.SlogMiddleware(logger))
	e.Use(middleware.ProblemMiddleware)
	e.Use(echom.CORSWithConfig(echom.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodOptions, http.MethodPut},
		AllowHeaders: []string{"*"},
	}))

	slog.SetDefault(logger)

	hh := NewHealthHandler(version, buildTime, commitHash)
	e.GET("/health", hh.Healthcheck)

	apidocs.APIDocsManager(e)

	return e
}
