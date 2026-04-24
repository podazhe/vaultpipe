package vault

import (
	"errors"
	"strings"
)

// NamespaceManager manages Vault namespace prefixing for secret paths.
type NamespaceManager struct {
	root      string
	namespace string
}

// NewNamespaceManager creates a NamespaceManager with the given root and active namespace.
// root is the base Vault namespace (e.g. "admin"), namespace is the sub-namespace.
func NewNamespaceManager(root, namespace string) (*NamespaceManager, error) {
	root = strings.Trim(root, "/")
	namespace = strings.Trim(namespace, "/")
	if root == "" {
		return nil, errors.New("namespace: root must not be empty")
	}
	return &NamespaceManager{root: root, namespace: namespace}, nil
}

// FullNamespace returns the fully-qualified namespace path.
func (nm *NamespaceManager) FullNamespace() string {
	if nm.namespace == "" {
		return nm.root
	}
	return nm.root + "/" + nm.namespace
}

// QualifyPath prepends the full namespace to the given secret path.
func (nm *NamespaceManager) QualifyPath(path string) string {
	path = strings.TrimPrefix(path, "/")
	return nm.FullNamespace() + "/" + path
}

// StripNamespace removes the full namespace prefix from a qualified path.
// Returns the original path unchanged if the prefix is not present.
func (nm *NamespaceManager) StripNamespace(qualifiedPath string) string {
	prefix := nm.FullNamespace() + "/"
	if strings.HasPrefix(qualifiedPath, prefix) {
		return qualifiedPath[len(prefix):]
	}
	return qualifiedPath
}

// QualifySecrets returns a new map with all keys qualified by the full namespace path.
func (nm *NamespaceManager) QualifySecrets(secrets map[string]string) map[string]string {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		out[nm.QualifyPath(k)] = v
	}
	return out
}

// StripSecrets returns a new map with the full namespace prefix removed from all keys.
func (nm *NamespaceManager) StripSecrets(secrets map[string]string) map[string]string {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		out[nm.StripNamespace(k)] = v
	}
	return out
}
