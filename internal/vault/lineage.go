package vault

import (
	"fmt"
	"sort"
	"time"
)

// LineageEntry records a single mutation event for a secret path.
type LineageEntry struct {
	Timestamp time.Time         `json:"timestamp"`
	Path      string            `json:"path"`
	Operation string            `json:"operation"` // "write", "delete", "promote", "rollback"
	Keys      []string          `json:"keys"`
	Meta      map[string]string `json:"meta,omitempty"`
}

// LineageTracker maintains an ordered log of secret mutation events
// so that callers can audit what changed, when, and why.
type LineageTracker struct {
	entries []LineageEntry
	maxSize int
}

// NewLineageTracker returns a LineageTracker that retains at most maxSize
// entries. A maxSize of 0 means unlimited.
func NewLineageTracker(maxSize int) *LineageTracker {
	return &LineageTracker{maxSize: maxSize}
}

// Record appends a new lineage entry. If the tracker is at capacity the
// oldest entry is evicted first.
func (lt *LineageTracker) Record(path, operation string, secrets map[string]string, meta map[string]string) {
	keys := make([]string, 0, len(secrets))
	for k := range secrets {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	entry := LineageEntry{
		Timestamp: time.Now().UTC(),
		Path:      path,
		Operation: operation,
		Keys:      keys,
		Meta:      meta,
	}

	if lt.maxSize > 0 && len(lt.entries) >= lt.maxSize {
		lt.entries = lt.entries[1:]
	}
	lt.entries = append(lt.entries, entry)
}

// History returns all recorded entries in chronological order.
func (lt *LineageTracker) History() []LineageEntry {
	out := make([]LineageEntry, len(lt.entries))
	copy(out, lt.entries)
	return out
}

// FilterByPath returns entries whose path matches the given value.
func (lt *LineageTracker) FilterByPath(path string) []LineageEntry {
	var out []LineageEntry
	for _, e := range lt.entries {
		if e.Path == path {
			out = append(out, e)
		}
	}
	return out
}

// FilterByOperation returns entries matching the given operation type.
func (lt *LineageTracker) FilterByOperation(op string) []LineageEntry {
	var out []LineageEntry
	for _, e := range lt.entries {
		if e.Operation == op {
			out = append(out, e)
		}
	}
	return out
}

// Summary returns a human-readable string describing all recorded events.
func (lt *LineageTracker) Summary() string {
	if len(lt.entries) == 0 {
		return "no lineage entries recorded"
	}
	out := fmt.Sprintf("%d lineage event(s):\n", len(lt.entries))
	for _, e := range lt.entries {
		out += fmt.Sprintf("  [%s] %s @ %s — %d key(s)\n",
			e.Timestamp.Format(time.RFC3339),
			e.Operation,
			e.Path,
			len(e.Keys),
		)
	}
	return out
}

// Clear removes all recorded entries.
func (lt *LineageTracker) Clear() {
	lt.entries = nil
}
