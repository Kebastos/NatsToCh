package log

import (
	"fmt"
	"log/slog"
	"os"
)

type Log struct {
	log *slog.Logger
}

func MustConfig() *Log {
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	return &Log{
		log: log,
	}
}

func (l *Log) SetDebug(d bool) {
	if d {
		l.log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	} else {
		l.log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
}

func (l *Log) Debugf(format string, args ...interface{}) {
	l.log.Debug(fmt.Sprintf(format, args...))
}

// Infof prints info message according to a format
func (l *Log) Infof(format string, args ...interface{}) {
	l.log.Info(fmt.Sprintf(format, args...))
}

// Errorf prints warning message according to a format
func (l *Log) Errorf(format string, args ...interface{}) {
	l.log.Error(fmt.Sprintf(format, args...))
}

// Fatalf prints fatal message according to a format and exits program
func (l *Log) Fatalf(format string, args ...interface{}) {
	l.log.Error(fmt.Sprintf(format, args...))
	os.Exit(1)
}
