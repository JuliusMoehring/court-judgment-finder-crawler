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

func (l *FileLogger) Debug(args ...interface{}) {
	l.logger.Debug(args...)
}

func (l *FileLogger) Debugf(format string, args ...interface{}) {
	l.logger.Debugf(format, args...)
}

func (l *FileLogger) Info(args ...interface{}) {
	l.logger.Info(args...)
}

func (l *FileLogger) Infof(format string, args ...interface{}) {
	l.logger.Infof(format, args...)
}

func (l *FileLogger) Warn(args ...interface{}) {
	l.logger.Warn(args...)
}

func (l *FileLogger) Warnf(format string, args ...interface{}) {
	l.logger.Warnf(format, args...)
}

func (l *FileLogger) Error(args ...interface{}) {
	l.logger.Error(args...)
}

func (l *FileLogger) Errorf(format string, args ...interface{}) {
	l.logger.Errorf(format, args...)
}

func (l *FileLogger) Fatal(args ...interface{}) {
	l.logger.Fatal(args...)
}

func (l *FileLogger) Fatalf(format string, args ...interface{}) {
	l.logger.Fatalf(format, args...)
}
