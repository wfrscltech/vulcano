//go:build !windows

package service

import (
	"log/slog"
)

func RunAsService(name string, log *slog.Logger, run func()) error {
	log.Info("Ejecutando en modo foreground (systemd compatible)", "name", name)
	run()
	return nil
}
