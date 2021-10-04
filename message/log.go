package message

import (
	"github.com/ThreeDotsLabs/watermill"
	"github.com/sirupsen/logrus"
)

// WatermillLogrusLogger is a wrapper of a Logrus logger implementing watermill.LoggerAdapter interface.
type WatermillLogrusLogger struct {
	*logrus.Logger
}

// NewWatermillLogrusLogger creates a new WatermillLogrusLogger.
func NewWatermillLogrusLogger(logger *logrus.Logger) *WatermillLogrusLogger {
	return &WatermillLogrusLogger{Logger: logger}
}

func (l WatermillLogrusLogger) Error(msg string, err error, fields watermill.LogFields) {
	l.Logger.WithError(err).WithFields(logrus.Fields(fields)).Error(msg)
}

func (l WatermillLogrusLogger) Info(msg string, fields watermill.LogFields) {
	l.Logger.WithFields(logrus.Fields(fields)).Info(msg)
}

func (l WatermillLogrusLogger) Debug(msg string, fields watermill.LogFields) {
	l.Logger.WithFields(logrus.Fields(fields)).Debug(msg)
}

func (l WatermillLogrusLogger) Trace(msg string, fields watermill.LogFields) {
	l.Logger.WithFields(logrus.Fields(fields)).Trace(msg)
}

func (l WatermillLogrusLogger) With(fields watermill.LogFields) watermill.LoggerAdapter {
	return &WatermillLogrusLogger{l.Logger.WithFields(logrus.Fields(fields)).Logger}
}
