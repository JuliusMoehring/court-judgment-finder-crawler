package logger

import (
	"os"

	log "github.com/sirupsen/logrus"
)

type FileLogger struct {
	logfile *os.File
	logger  *log.Logger
}

func NewFileLogger(filename string) Logger {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)

	if err != nil {
		panic("Could not open log file")
	}

	logger := log.New()

	logger.SetFormatter(&log.JSONFormatter{})
	logger.SetOutput(f)
	logger.SetLevel(log.DebugLevel)

	return &FileLogger{
		logfile: f,
		logger:  logger,
	}
}

func (l *FileLogger) Close() {
	l.logfile.Close()
}

func (l *FileLogger) Debug(domain string, args ...interface{}) {
	l.logger.Debug(args...)
}

func (l *FileLogger) Debugf(domain string, format string, args ...interface{}) {
	l.logger.Debugf(format, args...)
}

func (l *FileLogger) Info(domain string, args ...interface{}) {
	l.logger.Info(args...)
}

func (l *FileLogger) Infof(domain string, format string, args ...interface{}) {
	l.logger.Infof(format, args...)
}

func (l *FileLogger) Warn(domain string, args ...interface{}) {
	l.logger.Warn(args...)
}

func (l *FileLogger) Warnf(domain string, format string, args ...interface{}) {
	l.logger.Warnf(format, args...)
}

func (l *FileLogger) Error(domain string, args ...interface{}) {
	l.logger.Error(args...)
}

func (l *FileLogger) Errorf(domain string, format string, args ...interface{}) {
	l.logger.Errorf(format, args...)
}

func (l *FileLogger) Fatal(domain string, args ...interface{}) {
	l.logger.Fatal(args...)
}

func (l *FileLogger) Fatalf(domain string, format string, args ...interface{}) {
	l.logger.Fatalf(format, args...)
}
