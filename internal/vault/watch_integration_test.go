//go:build integration
// +build integration

package vault

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/vault/api"
)

func TestWatcher_Integration(t *testing.T) {
	addr := os.Getenv("VAULT_ADDR")
	token := os.Getenv("VAULT_TOKEN")
	path := os.Getenv("VAULT_WATCH_PATH")
	if addr == "" || token == "" || path == "" {
		t.Skip("VAULT_ADDR, VAULT_TOKEN, VAULT_WATCH_PATH required")
	}

	cfg := api.DefaultConfig()
	cfg.Address = addr
	client, err := api.NewClient(cfg)
	if err != nil {
		t.Fatalf("client: %v", err)
	}
	client.SetToken(token)

	watcher := NewWatcher(client, []string{path}, 2*time.Second)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	watcher.Start(ctx)
	defer watcher.Stop()

	select {
	case ev := <-watcher.Events():
		if ev.Err != nil {
			t.Fatalf("unexpected error: %v", ev.Err)
		}
		t.Logf("received event for %s with %d keys", ev.Path, len(ev.Data))
	case <-ctx.Done():
		t.Log("no change detected within timeout (expected if secret unchanged)")
	}
}
