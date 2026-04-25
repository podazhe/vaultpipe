package vault

import (
	"errors"
	"fmt"
)

// ImmutableSecrets wraps a secret map and prevents modification of locked keys.
type ImmutableSecrets struct {
	secrets map[string]string
	locked  map[string]bool
}

// NewImmutableSecrets creates an ImmutableSecrets instance from the provided map.
// The input map is copied to avoid external mutation.
func NewImmutableSecrets(secrets map[string]string) *ImmutableSecrets {
	s := make(map[string]string, len(secrets))
	for k, v := range secrets {
		s[k] = v
	}
	return &ImmutableSecrets{
		secrets: s,
		locked:  make(map[string]bool),
	}
}

// Lock marks one or more keys as immutable. Subsequent Set or Delete calls on
// locked keys will return an error.
func (im *ImmutableSecrets) Lock(keys ...string) {
	for _, k := range keys {
		im.locked[k] = true
	}
}

// IsLocked reports whether the given key is locked.
func (im *ImmutableSecrets) IsLocked(key string) bool {
	return im.locked[key]
}

// Get returns the value for key and whether it was found.
func (im *ImmutableSecrets) Get(key string) (string, bool) {
	v, ok := im.secrets[key]
	return v, ok
}

// Set updates or inserts a key. Returns an error if the key is locked.
func (im *ImmutableSecrets) Set(key, value string) error {
	if im.locked[key] {
		return fmt.Errorf("immutable: key %q is locked and cannot be modified", key)
	}
	im.secrets[key] = value
	return nil
}

// Delete removes a key. Returns an error if the key is locked.
func (im *ImmutableSecrets) Delete(key string) error {
	if im.locked[key] {
		return fmt.Errorf("immutable: key %q is locked and cannot be deleted", key)
	}
	delete(im.secrets, key)
	return nil
}

// Snapshot returns a plain copy of the underlying secret map.
func (im *ImmutableSecrets) Snapshot() map[string]string {
	out := make(map[string]string, len(im.secrets))
	for k, v := range im.secrets {
		out[k] = v
	}
	return out
}

// LockedKeys returns a sorted list of all locked keys.
func (im *ImmutableSecrets) LockedKeys() []string {
	keys := make([]string, 0, len(im.locked))
	for k := range im.locked {
		keys = append(keys, k)
	}
	return keys
}

// ErrNoLockedKeys is returned when an operation requires at least one locked key.
var ErrNoLockedKeys = errors.New("immutable: no keys are locked")
