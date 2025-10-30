package middleware

import (
	"log/slog"
	"time"

	"github.com/labstack/echo/v4"
)

// SlogMiddleware reemplaza al middleware.Logger() de Echo
func SlogMiddleware(logger *slog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			// Ejecutar el siguiente handler
			err := next(c)

			// Información básica de la request
			req := c.Request()
			res := c.Response()

			duration := time.Since(start)
			status := res.Status
			method := req.Method
			path := req.URL.Path
			ip := c.RealIP()

			// Nivel según código HTTP
			level := slog.LevelInfo
			if status >= 500 {
				level = slog.LevelError
			} else if status >= 400 {
				level = slog.LevelWarn
			}

			// Registrar usando slog
			logger.LogAttrs(
				c.Request().Context(),
				level,
				"HTTP request",
				slog.Int("status", status),
				slog.String("method", method),
				slog.String("path", path),
				slog.String("ip", ip),
				slog.String("latency", duration.String()),
			)

			return err
		}
	}
}
