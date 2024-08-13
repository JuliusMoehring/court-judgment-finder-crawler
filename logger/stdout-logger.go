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

func (l *StdOutLogger) Debug(domain string, args ...interface{}) {
	l.logger.WithFields(log.Fields{
		"domain": domain,
	}).Debug(args...)
}

func (l *StdOutLogger) Debugf(domain string, format string, args ...interface{}) {
	l.logger.WithFields(log.Fields{
		"domain": domain,
	}).Debugf(format, args...)
}

func (l *StdOutLogger) Info(domain string, args ...interface{}) {
	l.logger.WithFields(log.Fields{
		"domain": domain,
	}).Info(args...)
}

func (l *StdOutLogger) Infof(domain string, format string, args ...interface{}) {
	l.logger.WithFields(log.Fields{
		"domain": domain,
	}).Infof(format, args...)
}

func (l *StdOutLogger) Warn(domain string, args ...interface{}) {
	l.logger.WithFields(log.Fields{
		"domain": domain,
	}).Warn(args...)
}

func (l *StdOutLogger) Warnf(domain string, format string, args ...interface{}) {
	l.logger.WithFields(log.Fields{
		"domain": domain,
	}).Warnf(format, args...)
}

func (l *StdOutLogger) Error(domain string, args ...interface{}) {
	l.logger.WithFields(log.Fields{
		"domain": domain,
	}).Error(args...)
}

func (l *StdOutLogger) Errorf(domain string, format string, args ...interface{}) {
	l.logger.WithFields(log.Fields{
		"domain": domain,
	}).Errorf(format, args...)
}

func (l *StdOutLogger) Fatal(domain string, args ...interface{}) {
	l.logger.WithFields(log.Fields{
		"domain": domain,
	}).Fatal(args...)
}

func (l *StdOutLogger) Fatalf(domain string, format string, args ...interface{}) {
	l.logger.WithFields(log.Fields{
		"domain": domain,
	}).Fatalf(format, args...)
}
