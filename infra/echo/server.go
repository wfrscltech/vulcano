package echo

import (
	"context"
	"log/slog"

	"github.com/labstack/echo/v4"
)

type EchoServer struct {
	App  *echo.Echo
	Addr string
	Log  *slog.Logger
}

func (s EchoServer) Start() error {
	s.Log.Info("Servidor HTTP iniciado", "addr", s.Addr)
	return s.App.Start(s.Addr)
}

func (s EchoServer) Shutdown(ctx context.Context) error {
	s.Log.Info("Apagando servidor HTTP...")
	return s.App.Shutdown(ctx)
}
