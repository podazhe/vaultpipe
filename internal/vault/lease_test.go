package vault

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	vaultapi "github.com/hashicorp/vault/api"
)

func TestLeaseRenewer_NonRenewable(t *testing.T) {
	secret := &vaultapi.Secret{
		Renewable:     false,
		LeaseID:       "",
		LeaseDuration: 60,
	}
	client, err := vaultapi.NewClient(vaultapi.DefaultConfig())
	if err != nil {
		t.Fatal(err)
	}
	r := NewLeaseRenewer(client, secret)
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	r.Start(ctx)
	<-ctx.Done()
	// doneCh should be closed without blocking
	select {
	case <-r.doneCh:
		// ok — goroutine exited immediately
	case <-time.After(time.Second):
		t.Fatal("expected renewer to exit quickly for non-renewable secret")
	}
}

func TestLeaseRenewer_Stop(t *testing.T) {
	var renewCalls atomic.Int32

	srv := mockVaultServer(t, func(path string) map[string]interface{} {
		renewCalls.Add(1)
		return map[string]interface{}{}
	})

	cfg := vaultapi.DefaultConfig()
	cfg.Address = srv.URL
	client, err := vaultapi.NewClient(cfg)
	if err != nil {
		t.Fatal(err)
	}

	secret := &vaultapi.Secret{
		Renewable:     true,
		LeaseID:       "test/lease/123",
		LeaseDuration: 2, // 2s → renew every 1s
	}

	r := NewLeaseRenewer(client, secret)
	r.Start(context.Background())
	time.Sleep(150 * time.Millisecond)
	r.Stop()

	// doneCh must be closed after Stop
	select {
	case <-r.doneCh:
	case <-time.After(time.Second):
		t.Fatal("Stop() did not terminate the renewer in time")
	}
}
