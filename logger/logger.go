package logger

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"gopkg.in/natefinch/lumberjack.v2"
)

var Log *slog.Logger

func Init(level slog.Level, dir string) {
	_ = os.MkdirAll(dir, 0755)
	logFile := filepath.Join(dir, "vulcano.log")

	rotator := &lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    10,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   true,
	}

	mw := io.MultiWriter(os.Stdout, rotator)
	handler := slog.NewJSONHandler(mw, &slog.HandlerOptions{Level: level})
	Log = slog.New(handler)
}
