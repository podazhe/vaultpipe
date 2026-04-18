package vault

import (
	"context"
	"fmt"
	"time"

	vaultapi "github.com/hashicorp/vault/api"
)

// RotateConfig holds configuration for secret rotation.
type RotateConfig struct {
	Paths    []string
	Interval time.Duration
	OnRotate func(path string, data map[string]interface{})
}

// Rotator periodically re-fetches secrets and invokes a callback on change.
type Rotator struct {
	client *vaultapi.Client
	cache  *SecretCache
	cfg    RotateConfig
	stop   chan struct{}
}

// NewRotator creates a Rotator backed by the given client and cache.
func NewRotator(client *vaultapi.Client, cache *SecretCache, cfg RotateConfig) *Rotator {
	return &Rotator{
		client: client,
		cache:  cache,
		cfg:    cfg,
		stop:   make(chan struct{}),
	}
}

// Start begins the rotation loop; it blocks until Stop is called.
func (r *Rotator) Start(ctx context.Context) {
	ticker := time.NewTicker(r.cfg.Interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			r.poll(ctx)
		case <-r.stop:
			return
		case <-ctx.Done():
			return
		}
	}
}

// Stop signals the rotation loop to exit.
func (r *Rotator) Stop() { close(r.stop) }

func (r *Rotator) poll(ctx context.Context) {
	for _, path := range r.cfg.Paths {
		secret, err := r.client.Logical().ReadWithContext(ctx, kvv2Path(path))
		if err != nil || secret == nil {
			continue
		}
		data, err := extractData(secret)
		if err != nil {
			continue
		}
		key := fmt.Sprintf("rotate:%s", path)
		if _, hit := r.cache.Get(key); !hit {
			r.cache.Set(key, data, r.cfg.Interval*2)
			if r.cfg.OnRotate != nil {
				r.cfg.OnRotate(path, data)
			}
			continue
		}
		r.cache.Invalidate(key)
		r.cache.Set(key, data, r.cfg.Interval*2)
		if r.cfg.OnRotate != nil {
			r.cfg.OnRotate(path, data)
		}
	}
}
