package logger

import (
	"os"

	log "github.com/sirupsen/logrus"
)

type StdOutLogger struct {
	logger *log.Logger
}

func NewStdOutLogger() Logger {

	logger := log.New()

	logger.SetFormatter(&log.JSONFormatter{})
	logger.SetOutput(os.Stdout)
	logger.SetLevel(log.DebugLevel)

	return &StdOutLogger{
		logger: logger,
	}
}

func (l *StdOutLogger) Close() {
	// Nothing to do here
}

func (l *StdOutLogger) Debug(args ...interface{}) {
	l.logger.Debug(args...)
}

func (l *StdOutLogger) Debugf(format string, args ...interface{}) {
	l.logger.Debugf(format, args...)
}

func (l *StdOutLogger) Info(args ...interface{}) {
	l.logger.Info(args...)
}

func (l *StdOutLogger) Infof(format string, args ...interface{}) {
	l.logger.Infof(format, args...)
}

func (l *StdOutLogger) Warn(args ...interface{}) {
	l.logger.Warn(args...)
}

func (l *StdOutLogger) Warnf(format string, args ...interface{}) {
	l.logger.Warnf(format, args...)
}

func (l *StdOutLogger) Error(args ...interface{}) {
	l.logger.Error(args...)
}

func (l *StdOutLogger) Errorf(format string, args ...interface{}) {
	l.logger.Errorf(format, args...)
}

func (l *StdOutLogger) Fatal(args ...interface{}) {
	l.logger.Fatal(args...)
}

func (l *StdOutLogger) Fatalf(format string, args ...interface{}) {
	l.logger.Fatalf(format, args...)
}
