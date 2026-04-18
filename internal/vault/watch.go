package vault

import (
	"context"
	"time"

	"github.com/hashicorp/vault/api"
)

// WatchEvent is emitted when a secret changes or an error occurs.
type WatchEvent struct {
	Path   string
	Data   map[string]string
	Err    error
}

// Watcher polls a set of Vault paths and emits events when values change.
type Watcher struct {
	client   *api.Client
	paths    []string
	interval time.Duration
	events   chan WatchEvent
	stop     chan struct{}
}

// NewWatcher creates a Watcher that polls the given paths every interval.
func NewWatcher(client *api.Client, paths []string, interval time.Duration) *Watcher {
	return &Watcher{
		client:   client,
		paths:    paths,
		interval: interval,
		events:   make(chan WatchEvent, len(paths)),
		stop:     make(chan struct{}),
	}
}

// Events returns the read-only event channel.
func (w *Watcher) Events() <-chan WatchEvent {
	return w.events
}

// Start begins polling in a background goroutine.
func (w *Watcher) Start(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(w.interval)
		defer ticker.Stop()
		prev := make(map[string]map[string]string)
		for {
			select {
			case <-w.stop:
				return
			case <-ctx.Done():
				return
			case <-ticker.C:
				for _, p := range w.paths {
					data, err := ReadSecrets(w.client, p)
					if err != nil {
						w.events <- WatchEvent{Path: p, Err: err}
						continue
					}
					if !mapsEqual(prev[p], data) {
						prev[p] = data
						w.events <- WatchEvent{Path: p, Data: data}
					}
				}
			}
		}
	}()
}

// Stop halts the watcher.
func (w *Watcher) Stop() {
	close(w.stop)
}

func mapsEqual(a, b map[string]string) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if b[k] != v {
			return false
		}
	}
	return true
}
