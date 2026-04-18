package vault

import (
	"sync"
	"time"
)

// CacheEntry holds a cached secret and its expiry.
type CacheEntry struct {
	Data      map[string]string
	FetchedAt time.Time
	TTL       time.Duration
}

// IsExpired returns true if the cache entry has passed its TTL.
func (e *CacheEntry) IsExpired() bool {
	if e.TTL <= 0 {
		return false
	}
	return time.Since(e.FetchedAt) > e.TTL
}

// SecretCache is a thread-safe in-memory cache for Vault secrets.
type SecretCache struct {
	mu      sync.RWMutex
	entries map[string]*CacheEntry
	defaultTTL time.Duration
}

// NewSecretCache creates a new SecretCache with the given default TTL.
// Pass 0 to disable expiry.
func NewSecretCache(defaultTTL time.Duration) *SecretCache {
	return &SecretCache{
		entries:    make(map[string]*CacheEntry),
		defaultTTL: defaultTTL,
	}
}

// Get retrieves a cached secret by path. Returns nil if missing or expired.
func (c *SecretCache) Get(path string) map[string]string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	e, ok := c.entries[path]
	if !ok || e.IsExpired() {
		return nil
	}
	return e.Data
}

// Set stores a secret in the cache under the given path.
func (c *SecretCache) Set(path string, data map[string]string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[path] = &CacheEntry{
		Data:      data,
		FetchedAt: time.Now(),
		TTL:       c.defaultTTL,
	}
}

// Invalidate removes a specific path from the cache.
func (c *SecretCache) Invalidate(path string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.entries, path)
}

// Flush clears all entries from the cache.
func (c *SecretCache) Flush() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries = make(map[string]*CacheEntry)
}
