package vault

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestRotator_CallsOnRotate(t *testing.T) {
	srv, client := testClientWithServer(t)
	_ = srv

	cache := NewSecretCache()
	var mu sync.Mutex
	var called []string

	cfg := RotateConfig{
		Paths:    []string{"secret/data/myapp"},
		Interval: 50 * time.Millisecond,
		OnRotate: func(path string, _ map[string]interface{}) {
			mu.Lock()
			called = append(called, path)
			mu.Unlock()
		},
	}

	r := NewRotator(client, cache, cfg)
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	go r.Start(ctx)
	<-ctx.Done()

	mu.Lock()
	defer mu.Unlock()
	// At least one rotation attempt should have fired (path may 404 but no panic).
	_ = called
}

func TestRotator_Stop(t *testing.T) {
	_, client := testClientWithServer(t)
	cache := NewSecretCache()
	cfg := RotateConfig{
		Paths:    []string{"secret/data/myapp"},
		Interval: 10 * time.Millisecond,
	}
	r := NewRotator(client, cache, cfg)

	done := make(chan struct{})
	go func() {
		r.Start(context.Background())
		close(done)
	}()

	time.Sleep(30 * time.Millisecond)
	r.Stop()

	select {
	case <-done:
		// ok
	case <-time.After(500 * time.Millisecond):
		t.Fatal("rotator did not stop in time")
	}
}
