package vault

import (
	"fmt"
	"sync"
	"time"
)

// RollbackEntry holds a versioned snapshot of secrets.
type RollbackEntry struct {
	Version   int
	Timestamp time.Time
	Secrets   map[string]string
}

// RollbackManager maintains a history of secret states for rollback.
type RollbackManager struct {
	mu      sync.Mutex
	history []RollbackEntry
	maxSize int
}

// NewRollbackManager creates a RollbackManager with a capped history size.
func NewRollbackManager(maxSize int) *RollbackManager {
	if maxSize <= 0 {
		maxSize = 10
	}
	return &RollbackManager{maxSize: maxSize}
}

// Push records a new version of secrets.
func (r *RollbackManager) Push(secrets map[string]string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	copy := make(map[string]string, len(secrets))
	for k, v := range secrets {
		copy[k] = v
	}
	version := 1
	if len(r.history) > 0 {
		version = r.history[len(r.history)-1].Version + 1
	}
	r.history = append(r.history, RollbackEntry{
		Version:   version,
		Timestamp: time.Now(),
		Secrets:   copy,
	})
	if len(r.history) > r.maxSize {
		r.history = r.history[len(r.history)-r.maxSize:]
	}
}

// Rollback returns secrets at the given version number.
func (r *RollbackManager) Rollback(version int) (map[string]string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, entry := range r.history {
		if entry.Version == version {
			copy := make(map[string]string, len(entry.Secrets))
			for k, v := range entry.Secrets {
				copy[k] = v
			}
			return copy, nil
		}
	}
	return nil, fmt.Errorf("rollback: version %d not found", version)
}

// History returns all stored entries.
func (r *RollbackManager) History() []RollbackEntry {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make([]RollbackEntry, len(r.history))
	copy(out, r.history)
	return out
}
