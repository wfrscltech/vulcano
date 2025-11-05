package logger

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/natefinch/lumberjack.v2"
)

var Log *slog.Logger

func Init(level slog.Level, logname, dir string) {
	_ = os.MkdirAll(dir, 0755)
	var out io.Writer = os.Stdout

	if strings.HasPrefix(dir, "dir:") {
		logFile := filepath.Join(dir[4:], logname+".log")
		rotator := &lumberjack.Logger{
			Filename:   logFile,
			MaxSize:    10,
			MaxBackups: 5,
			MaxAge:     30,
			Compress:   true,
		}

		out = io.MultiWriter(os.Stdout, rotator)
	}
	handler := slog.NewJSONHandler(out, &slog.HandlerOptions{Level: level})
	Log = slog.New(handler)
}
