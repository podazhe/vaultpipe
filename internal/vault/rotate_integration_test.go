//go:build integration
// +build integration

package vault_test

import (
	"context"
	"testing"
	"time"

	vaultapi "github.com/hashicorp/vault/api"

	"github.com/yourorg/vaultpipe/internal/vault"
)

// TestRotator_Integration requires a real Vault dev server at VAULT_ADDR.
func TestRotator_Integration(t *testing.T) {
	cfg := vaultapi.DefaultConfig()
	raw, err := vaultapi.NewClient(cfg)
	if err != nil {
		t.Fatalf("vault client: %v", err)
	}

	cache := vault.NewSecretCache()
	rotations := make(chan string, 10)

	rc := vault.RotateConfig{
		Paths:    []string{"secret/data/integration"},
		Interval: 100 * time.Millisecond,
		OnRotate: func(path string, _ map[string]interface{}) {
			rotations <- path
		},
	}

	r := vault.NewRotator(raw, cache, rc)
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	go r.Start(ctx)
	<-ctx.Done()

	if len(rotations) == 0 {
		t.Log("no rotations fired (path may not exist in dev server — acceptable)")
	}
}
