package vault

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"
)

// Checkpoint represents a named point-in-time snapshot of secrets.
type Checkpoint struct {
	Name      string            `json:"name"`
	CreatedAt time.Time         `json:"created_at"`
	Secrets   map[string]string `json:"secrets"`
}

// CheckpointManager manages named checkpoints persisted to disk.
type CheckpointManager struct {
	dir string
}

// NewCheckpointManager creates a CheckpointManager rooted at dir.
func NewCheckpointManager(dir string) (*CheckpointManager, error) {
	if dir == "" {
		return nil, errors.New("checkpoint: directory must not be empty")
	}
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, fmt.Errorf("checkpoint: create dir: %w", err)
	}
	return &CheckpointManager{dir: dir}, nil
}

// Save persists secrets under the given name.
func (m *CheckpointManager) Save(name string, secrets map[string]string) error {
	if name == "" {
		return errors.New("checkpoint: name must not be empty")
	}
	cp := Checkpoint{
		Name:      name,
		CreatedAt: time.Now().UTC(),
		Secrets:   secrets,
	}
	data, err := json.MarshalIndent(cp, "", "  ")
	if err != nil {
		return fmt.Errorf("checkpoint: marshal: %w", err)
	}
	return os.WriteFile(m.filePath(name), data, 0600)
}

// Load retrieves a checkpoint by name.
func (m *CheckpointManager) Load(name string) (*Checkpoint, error) {
	if name == "" {
		return nil, errors.New("checkpoint: name must not be empty")
	}
	data, err := os.ReadFile(m.filePath(name))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("checkpoint %q not found", name)
		}
		return nil, fmt.Errorf("checkpoint: read: %w", err)
	}
	var cp Checkpoint
	if err := json.Unmarshal(data, &cp); err != nil {
		return nil, fmt.Errorf("checkpoint: unmarshal: %w", err)
	}
	return &cp, nil
}

// Delete removes a checkpoint by name.
func (m *CheckpointManager) Delete(name string) error {
	if name == "" {
		return errors.New("checkpoint: name must not be empty")
	}
	err := os.Remove(m.filePath(name))
	if errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("checkpoint %q not found", name)
	}
	return err
}

func (m *CheckpointManager) filePath(name string) string {
	return fmt.Sprintf("%s/%s.json", m.dir, name)
}
