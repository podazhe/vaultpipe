package vault

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

// AuditEvent represents a single recorded operation on a secret.
type AuditEvent struct {
	Timestamp time.Time
	Operation string // e.g. "read", "write", "delete", "rotate"
	Path      string
	Key       string
	Actor     string
	Success   bool
	Message   string
}

// AuditTrail records and queries a history of secret operations.
type AuditTrail struct {
	events []AuditEvent
	maxEvents int
}

// AuditTrailOption configures an AuditTrail.
type AuditTrailOption func(*AuditTrail)

// WithMaxEvents caps the number of events retained (oldest are dropped).
func WithMaxEvents(n int) AuditTrailOption {
	return func(a *AuditTrail) {
		if n > 0 {
			a.maxEvents = n
		}
	}
}

// NewAuditTrail creates an AuditTrail with the given options.
func NewAuditTrail(opts ...AuditTrailOption) *AuditTrail {
	at := &AuditTrail{
		maxEvents: 10_000,
	}
	for _, o := range opts {
		o(at)
	}
	return at
}

// Record appends a new event to the trail.
func (a *AuditTrail) Record(op, path, key, actor string, success bool, msg string) {
	event := AuditEvent{
		Timestamp: time.Now().UTC(),
		Operation: op,
		Path:      path,
		Key:       key,
		Actor:     actor,
		Success:   success,
		Message:   msg,
	}
	a.events = append(a.events, event)
	if len(a.events) > a.maxEvents {
		a.events = a.events[len(a.events)-a.maxEvents:]
	}
}

// All returns a copy of all recorded events in chronological order.
func (a *AuditTrail) All() []AuditEvent {
	out := make([]AuditEvent, len(a.events))
	copy(out, a.events)
	return out
}

// FilterByPath returns events whose Path matches the given prefix.
func (a *AuditTrail) FilterByPath(prefix string) []AuditEvent {
	var out []AuditEvent
	for _, e := range a.events {
		if strings.HasPrefix(e.Path, prefix) {
			out = append(out, e)
		}
	}
	return out
}

// FilterByActor returns events recorded for the given actor.
func (a *AuditTrail) FilterByActor(actor string) []AuditEvent {
	var out []AuditEvent
	for _, e := range a.events {
		if e.Actor == actor {
			out = append(out, e)
		}
	}
	return out
}

// FilterByOperation returns events matching the given operation type.
func (a *AuditTrail) FilterByOperation(op string) []AuditEvent {
	var out []AuditEvent
	for _, e := range a.events {
		if e.Operation == op {
			out = append(out, e)
		}
	}
	return out
}

// Summary returns a human-readable breakdown of event counts per operation.
func (a *AuditTrail) Summary() string {
	counts := make(map[string]int)
	for _, e := range a.events {
		counts[e.Operation]++
	}
	ops := make([]string, 0, len(counts))
	for op := range counts {
		ops = append(ops, op)
	}
	sort.Strings(ops)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("audit trail: %d event(s)\n", len(a.events)))
	for _, op := range ops {
		sb.WriteString(fmt.Sprintf("  %-12s %d\n", op, counts[op]))
	}
	return strings.TrimRight(sb.String(), "\n")
}

// Clear removes all recorded events.
func (a *AuditTrail) Clear() {
	a.events = nil
}
