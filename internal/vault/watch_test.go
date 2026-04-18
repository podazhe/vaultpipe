package vault

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

func TestWatcher_EmitsOnChange(t *testing.T) {
	var callCount int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := atomic.AddInt32(&callCount, 1)
		value := "first"
		if n > 1 {
			value = "second"
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"data":{"data":{"KEY":"` + value + `"},"metadata":{}}}`)) 
	}))
	defer server.Close()

	client := testClientWithServer(t, server)
	watcher := NewWatcher(client, []string{"secret/data/app"}, 50*time.Millisecond)
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	watcher.Start(ctx)
	defer watcher.Stop()

	var events []WatchEvent
	timeout := time.After(400 * time.Millisecond)
loop:
	for {
		select {
		case ev := <-watcher.Events():
			events = append(events, ev)
			if len(events) >= 2 {
				break loop
			}
		case <-timeout:
			break loop
		}
	}

	if len(events) < 2 {
		t.Fatalf("expected at least 2 events, got %d", len(events))
	}
	if events[0].Data["KEY"] != "first" {
		t.Errorf("expected first, got %s", events[0].Data["KEY"])
	}
	if events[1].Data["KEY"] != "second" {
		t.Errorf("expected second, got %s", events[1].Data["KEY"])
	}
}

func TestWatcher_Stop(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"data":{"data":{"K":"v"},"metadata":{}}}`)) 
	}))
	defer server.Close()
	client := testClientWithServer(t, server)
	watcher := NewWatcher(client, []string{"secret/data/app"}, 20*time.Millisecond)
	watcher.Start(context.Background())
	time.Sleep(60 * time.Millisecond)
	watcher.Stop()
	// Should not panic after stop
}

func TestMapsEqual(t *testing.T) {
	if !mapsEqual(map[string]string{"a": "1"}, map[string]string{"a": "1"}) {
		t.Error("expected equal")
	}
	if mapsEqual(map[string]string{"a": "1"}, map[string]string{"a": "2"}) {
		t.Error("expected not equal")
	}
}
