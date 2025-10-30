package service

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Runner interface {
	Start() error
	Shutdown(ctx context.Context) error
}

// RunGracefully inicia un servicio y lo mantiene en ejecución hasta recibir una señal de
// apagado.
func RunGracefully(log *slog.Logger, srv Runner) error {
	go func() {
		if err := srv.Start(); err != nil {
			log.Error("Error al iniciar servicio", "error", err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Info("Se recibe señal de apagado")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return srv.Shutdown(ctx)
}
