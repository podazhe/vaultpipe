package vault

import (
	"fmt"
	"sync"
	"time"
)

// PinnedSecret holds a secret value locked to a specific version.
type PinnedSecret struct {
	Path      string
	Version   int
	Data      map[string]string
	PinnedAt  time.Time
}

// PinManager tracks secrets pinned to specific versions.
type PinManager struct {
	mu   sync.RWMutex
	pins map[string]*PinnedSecret
}

// NewPinManager creates a new PinManager.
func NewPinManager() *PinManager {
	return &PinManager{
		pins: make(map[string]*PinnedSecret),
	}
}

// Pin records a secret path at a given version with its data.
func (pm *PinManager) Pin(path string, version int, data map[string]string) error {
	if path == "" {
		return fmt.Errorf("pin: path must not be empty")
	}
	if version < 1 {
		return fmt.Errorf("pin: version must be >= 1, got %d", version)
	}
	copy := make(map[string]string, len(data))
	for k, v := range data {
		copy[k] = v
	}
	pm.mu.Lock()
	defer pm.mu.Unlock()
	pm.pins[path] = &PinnedSecret{
		Path:     path,
		Version:  version,
		Data:     copy,
		PinnedAt: time.Now(),
	}
	return nil
}

// Get returns the pinned secret for a path, or an error if not pinned.
func (pm *PinManager) Get(path string) (*PinnedSecret, error) {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	p, ok := pm.pins[path]
	if !ok {
		return nil, fmt.Errorf("pin: no pin found for path %q", path)
	}
	return p, nil
}

// Unpin removes the pin for a path.
func (pm *PinManager) Unpin(path string) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	delete(pm.pins, path)
}

// List returns all currently pinned paths.
func (pm *PinManager) List() []string {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	paths := make([]string, 0, len(pm.pins))
	for p := range pm.pins {
		paths = append(paths, p)
	}
	return paths
}
