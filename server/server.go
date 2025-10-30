package server

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Interfaz base para servidores (Echo u otros frameworks).

type HTTPServer struct {
	Addr string
	Log  *slog.Logger
	Mux  http.Handler
}

func (s *HTTPServer) Start() error {
	s.Log.Info("Servidor HTTP iniciado", "addr", s.Addr)
	return http.ListenAndServe(s.Addr, s.Mux)
}

func (s *HTTPServer) Shutdown(ctx context.Context) error {
	s.Log.Info("Apagando servidor HTTP...")
	return nil
}

func (s *HTTPServer) RunGracefully() error {
	go func() {
		if err := s.Start(); err != nil && err != http.ErrServerClosed {
			s.Log.Error("Error en servidor", "error", err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	s.Log.Info("Recibida señal de apagado")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return s.Shutdown(ctx)
}
