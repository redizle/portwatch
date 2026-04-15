package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Event represents a port activity event to be logged
type Event struct {
	Timestamp time.Time `json:"timestamp"`
	Port      int       `json:"port"`
	Status    string    `json:"status"` // "open" or "closed"
	Host      string    `json:"host"`
}

// Logger handles writing port events to a log destination
type Logger struct {
	file   *os.File
	format string // "json" or "text"
}

// New creates a new Logger. If logPath is empty, logs go to stdout.
func New(logPath string, format string) (*Logger, error) {
	if format != "json" && format != "text" {
		format = "text"
	}

	var f *os.File
	var err error

	if logPath == "" {
		f = os.Stdout
	} else {
		f, err = os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return nil, fmt.Errorf("failed to open log file: %w", err)
		}
	}

	return &Logger{file: f, format: format}, nil
}

// Log writes a port event to the log destination
func (l *Logger) Log(e Event) error {
	if e.Timestamp.IsZero() {
		e.Timestamp = time.Now()
	}

	var line string

	if l.format == "json" {
		b, err := json.Marshal(e)
		if err != nil {
			return fmt.Errorf("failed to marshal event: %w", err)
		}
		line = string(b) + "\n"
	} else {
		line = fmt.Sprintf("%s  host=%s port=%d status=%s\n",
			e.Timestamp.Format(time.RFC3339), e.Host, e.Port, e.Status)
	}

	_, err := fmt.Fprint(l.file, line)
	return err
}

// Close closes the underlying log file if it's not stdout
func (l *Logger) Close() error {
	if l.file != nil && l.file != os.Stdout {
		return l.file.Close()
	}
	return nil
}
