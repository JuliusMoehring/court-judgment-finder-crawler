package logger

type Logger interface {
	Close()

	Debug(domain string, args ...interface{})
	Debugf(domain string, format string, args ...interface{})

	Info(domain string, args ...interface{})
	Infof(domain string, format string, args ...interface{})

	Warn(domain string, args ...interface{})
	Warnf(domain string, format string, args ...interface{})

	Error(domain string, args ...interface{})
	Errorf(domain string, format string, args ...interface{})

	Fatal(domain string, args ...interface{})
	Fatalf(domain string, format string, args ...interface{})
}
