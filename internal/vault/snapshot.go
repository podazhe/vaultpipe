package vault

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Snapshot holds a point-in-time capture of secrets at a given path.
type Snapshot struct {
	Path      string            `json:"path"`
	Data      map[string]string `json:"data"`
	CapturedAt time.Time        `json:"captured_at"`
}

// SnapshotManager saves and loads secret snapshots to/from disk.
type SnapshotManager struct {
	dir string
}

func NewSnapshotManager(dir string) *SnapshotManager {
	return &SnapshotManager{dir: dir}
}

func (m *SnapshotManager) Save(path string, data map[string]string) error {
	snap := Snapshot{
		Path:       path,
		Data:       data,
		CapturedAt: time.Now().UTC(),
	}
	b, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return fmt.Errorf("snapshot marshal: %w", err)
	}
	filePath := m.filePath(path)
	if err := os.MkdirAll(m.dir, 0700); err != nil {
		return fmt.Errorf("snapshot mkdir: %w", err)
	}
	return os.WriteFile(filePath, b, 0600)
}

func (m *SnapshotManager) Load(path string) (*Snapshot, error) {
	b, err := os.ReadFile(m.filePath(path))
	if err != nil {
		return nil, fmt.Errorf("snapshot read: %w", err)
	}
	var snap Snapshot
	if err := json.Unmarshal(b, &snap); err != nil {
		return nil, fmt.Errorf("snapshot unmarshal: %w", err)
	}
	return &snap, nil
}

func (m *SnapshotManager) filePath(secretPath string) string {
	safe := sanitiseKey(secretPath)
	return fmt.Sprintf("%s/%s.json", m.dir, safe)
}
