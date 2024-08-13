package logger

import (
	"bufio"
	"encoding/json"
	"os"
	"testing"
	"time"
)

const LOG_FILE = "test.log"

type Log struct {
	Level   string    `json:"level"`
	Message string    `json:"msg"`
	Time    time.Time `json:"time"`
}

func TestFileLogger(t *testing.T) {
	assertLogFileExists := func(t *testing.T) {
		t.Helper()
		_, err := os.Stat(LOG_FILE)
		if err != nil {
			t.Errorf("expected log file to exist, but it does not")
		}
	}

	assertLogContains := func(t *testing.T, level string, message string) {
		t.Helper()
		file, err := os.Open(LOG_FILE)
		if err != nil {
			t.Errorf("expected to be able to open log file, but could not")
		}
		defer file.Close()

		reader := bufio.NewReader(file)
		found := false

		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				if err.Error() == "EOF" {
					break
				}
				t.Errorf("expected to be able to read log file, but could not")
			}

			var log Log
			err = json.Unmarshal([]byte(line), &log)
			if err != nil {
				t.Errorf("expected to be able to unmarshal log line, but could not")
			}

			if log.Level == level && log.Message == message {
				found = true
				break
			}
		}

		if !found {
			t.Errorf("expected log file to contain log entry with level %s and message %s, but it did not", level, message)
		}
	}

	t.Cleanup(func() {
		os.Remove(LOG_FILE)
	})

	t.Run("Logs are written to log file", func(t *testing.T) {
		logger := NewFileLogger(LOG_FILE)

		logger.Debug("This is a debug message")
		logger.Info("This is an info message")
		logger.Warn("This is a warning message")
		logger.Error("This is an error message")

		logger.Close()

		assertLogFileExists(t)
		assertLogContains(t, "debug", "This is a debug message")
		assertLogContains(t, "info", "This is an info message")
		assertLogContains(t, "warning", "This is a warning message")
		assertLogContains(t, "error", "This is an error message")
	})
}
