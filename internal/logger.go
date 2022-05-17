package internal

import (
	"github.com/chatex-com/process-manager"
	log "github.com/sirupsen/logrus"
)

type ProcessLogger struct {
	Logger *log.Logger
}

func NewProcessLogger() *ProcessLogger {
	return &ProcessLogger{
		Logger: log.StandardLogger(),
	}
}

func (pl *ProcessLogger) Info(msg string, fields ...process.LogFields) {
	entry := log.NewEntry(pl.Logger)

	if len(fields) > 0 {
		entry = entry.WithFields(log.Fields(fields[0]))
	}

	entry.Info(msg)
}

func (pl *ProcessLogger) Error(msg string, err error, fields ...process.LogFields) {
	entry := log.NewEntry(pl.Logger)

	if len(fields) > 0 {
		entry = entry.WithFields(log.Fields(fields[0]))
	}

	entry.WithError(err).Error(msg)
}
