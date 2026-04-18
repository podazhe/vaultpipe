package audit

import (
	"encoding/json"
	"io"
	"os"
	"time"
)

// Level represents the severity of an audit event.
type Level string

const (
	LevelInfo  Level = "INFO"
	LevelWarn  Level = "WARN"
	LevelError Level = "ERROR"
)

// Event captures a single auditable action within vaultpipe.
type Event struct {
	Timestamp time.Time         `json:"timestamp"`
	Level     Level             `json:"level"`
	Action    string            `json:"action"`
	Path      string            `json:"path,omitempty"`
	Meta      map[string]string `json:"meta,omitempty"`
	Error     string            `json:"error,omitempty"`
}

// Logger writes structured audit events as JSON lines.
type Logger struct {
	out io.Writer
}

// NewLogger returns a Logger writing to w. Pass nil to use stderr.
func NewLogger(w io.Writer) *Logger {
	if w == nil {
		w = os.Stderr
	}
	return &Logger{out: w}
}

// Log emits an audit event.
func (l *Logger) Log(level Level, action, path string, meta map[string]string, err error) {
	e := Event{
		Timestamp: time.Now().UTC(),
		Level:     level,
		Action:    action,
		Path:      path,
		Meta:      meta,
	}
	if err != nil {
		e.Error = err.Error()
	}
	data, _ := json.Marshal(e)
	_, _ = l.out.Write(append(data, '\n'))
}

// Info is a convenience wrapper for LevelInfo events.
func (l *Logger) Info(action, path string, meta map[string]string) {
	l.Log(LevelInfo, action, path, meta, nil)
}

// Warn is a convenience wrapper for LevelWarn events.
func (l *Logger) Warn(action, path string, err error) {
	l.Log(LevelWarn, action, path, nil, err)
}

// Error is a convenience wrapper for LevelError events.
func (l *Logger) Error(action, path string, err error) {
	l.Log(LevelError, action, path, nil, err)
}
