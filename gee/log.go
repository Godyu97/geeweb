package gee

import (
	"golang.org/x/exp/slog"
	"os"
)

type Logger interface {
	Info(msg string, a ...any)
}

var geeLogger Logger = slog.New(slog.NewTextHandler(os.Stderr, nil))

func initLogger(logger Logger) {
	//todo logger init
	geeLogger = logger
}
