package vault

import (
	"fmt"
	"strings"
)

// Scope represents a named boundary for grouping secrets by environment or team.
type Scope struct {
	Name   string
	Prefix string
	Tags   map[string]string
}

// ScopeManager manages multiple named scopes and maps secret keys to them.
type ScopeManager struct {
	scopes map[string]*Scope
}

// NewScopeManager creates a new ScopeManager with no scopes registered.
func NewScopeManager() *ScopeManager {
	return &ScopeManager{
		scopes: make(map[string]*Scope),
	}
}

// Register adds a named scope with an optional key prefix and tags.
func (sm *ScopeManager) Register(name, prefix string, tags map[string]string) error {
	if name == "" {
		return fmt.Errorf("scope name must not be empty")
	}
	if _, exists := sm.scopes[name]; exists {
		return fmt.Errorf("scope %q already registered", name)
	}
	sm.scopes[name] = &Scope{
		Name:   name,
		Prefix: prefix,
		Tags:   tags,
	}
	return nil
}

// Resolve returns the scope that best matches the given secret key.
// Matching is based on the longest prefix that covers the key.
func (sm *ScopeManager) Resolve(key string) (*Scope, bool) {
	var best *Scope
	for _, s := range sm.scopes {
		if s.Prefix == "" {
			continue
		}
		if strings.HasPrefix(key, s.Prefix) {
			if best == nil || len(s.Prefix) > len(best.Prefix) {
				best = s
			}
		}
	}
	if best != nil {
		return best, true
	}
	return nil, false
}

// Partition groups a secrets map into per-scope buckets.
// Keys that match no scope are placed under the "_default" bucket.
func (sm *ScopeManager) Partition(secrets map[string]string) map[string]map[string]string {
	result := make(map[string]map[string]string)
	for k, v := range secrets {
		scope, ok := sm.Resolve(k)
		bucket := "_default"
		if ok {
			bucket = scope.Name
		}
		if result[bucket] == nil {
			result[bucket] = make(map[string]string)
		}
		result[bucket][k] = v
	}
	return result
}

// List returns all registered scope names.
func (sm *ScopeManager) List() []string {
	names := make([]string, 0, len(sm.scopes))
	for n := range sm.scopes {
		names = append(names, n)
	}
	return names
}
