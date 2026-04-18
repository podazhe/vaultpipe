package vault

import (
	"context"
	"log"
	"time"

	vaultapi "github.com/hashicorp/vault/api"
)

// LeaseRenewer watches a secret lease and renews it before expiry.
type LeaseRenewer struct {
	client  *vaultapi.Client
	secret  *vaultapi.Secret
	stopCh  chan struct{}
	doneCh  chan struct{}
}

// NewLeaseRenewer creates a LeaseRenewer for the given renewable secret.
func NewLeaseRenewer(client *vaultapi.Client, secret *vaultapi.Secret) *LeaseRenewer {
	return &LeaseRenewer{
		client: client,
		secret: secret,
		stopCh: make(chan struct{}),
		doneCh: make(chan struct{}),
	}
}

// Start begins the renewal loop in a goroutine.
// It renews the lease at half the LeaseDuration interval.
func (r *LeaseRenewer) Start(ctx context.Context) {
	go func() {
		defer close(r.doneCh)
		if !r.secret.Renewable || r.secret.LeaseID == "" {
			return
		}
		interval := time.Duration(r.secret.LeaseDuration/2) * time.Second
		if interval <= 0 {
			return
		}
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-r.stopCh:
				return
			case <-ticker.C:
				_, err := r.client.Sys().Renew(r.secret.LeaseID, r.secret.LeaseDuration)
				if err != nil {
					log.Printf("vaultpipe: failed to renew lease %s: %v", r.secret.LeaseID, err)
					return
				}
				log.Printf("vaultpipe: renewed lease %s", r.secret.LeaseID)
			}
		}
	}()
}

// Stop signals the renewal loop to exit and waits for it to finish.
func (r *LeaseRenewer) Stop() {
	close(r.stopCh)
	<-r.doneCh
}
