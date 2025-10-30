//go:build windows

package service

import (
	"log/slog"
	"os"

	"golang.org/x/sys/windows/svc"
)

func RunAsService(name string, log *slog.Logger, run func()) error {
	if !svc.IsWindowsService() {
		run()
		return nil
	}
	log.Info("Ejecutando como servicio de Windows", "name", name)
	return runService(name, run)
}

func runService(name string, run func()) error {
	return svc.Run(name, &program{run: run})
}

type program struct {
	run func()
}

func (p *program) Execute(args []string, r <-chan svc.ChangeRequest, s chan<- svc.Status) (bool, uint32) {
	const cmdsAccepted = svc.AcceptStop | svc.AcceptShutdown
	s <- svc.Status{State: svc.StartPending}
	go p.run()
	s <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}
	for c := range r {
		switch c.Cmd {
		case svc.Interrogate:
			s <- c.CurrentStatus
		case svc.Stop, svc.Shutdown:
			s <- svc.Status{State: svc.StopPending}
			os.Exit(0)
			return false, 0
		}
	}
	return false, 0
}
