package audit

import (
	"bytes"
	"encoding/json"
	"errors"
	"testing"
)

func newTestLogger() (*Logger, *bytes.Buffer) {
	var buf bytes.Buffer
	return NewLogger(&buf), &buf
}

func TestLogger_Info(t *testing.T) {
	l, buf := newTestLogger()
	l.Info("read_secret", "secret/data/app", map[string]string{"keys": "3"})

	var e Event
	if err := json.Unmarshal(buf.Bytes(), &e); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if e.Level != LevelInfo {
		t.Errorf("expected INFO, got %s", e.Level)
	}
	if e.Action != "read_secret" {
		t.Errorf("unexpected action: %s", e.Action)
	}
	if e.Path != "secret/data/app" {
		t.Errorf("unexpected path: %s", e.Path)
	}
	if e.Meta["keys"] != "3" {
		t.Errorf("meta not propagated")
	}
	if e.Error != "" {
		t.Errorf("expected no error field")
	}
}

func TestLogger_Error(t *testing.T) {
	l, buf := newTestLogger()
	sentinel := errors.New("permission denied")
	l.Error("read_secret", "secret/data/restricted", sentinel)

	var e Event
	if err := json.Unmarshal(buf.Bytes(), &e); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if e.Level != LevelError {
		t.Errorf("expected ERROR, got %s", e.Level)
	}
	if e.Error != "permission denied" {
		t.Errorf("unexpected error field: %s", e.Error)
	}
}

func TestLogger_Warn(t *testing.T) {
	l, buf := newTestLogger()
	l.Warn("lease_renew", "secret/data/db", errors.New("lease not renewable"))

	var e Event
	if err := json.Unmarshal(buf.Bytes(), &e); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if e.Level != LevelWarn {
		t.Errorf("expected WARN, got %s", e.Level)
	}
}

func TestLogger_TimestampPresent(t *testing.T) {
	l, buf := newTestLogger()
	l.Info("startup", "", nil)

	var e Event
	if err := json.Unmarshal(buf.Bytes(), &e); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if e.Timestamp.IsZero() {
		t.Error("timestamp should not be zero")
	}
}
