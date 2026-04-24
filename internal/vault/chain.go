package vault

import "fmt"

// ChainStep represents a single transformation step in a secret processing chain.
type ChainStep struct {
	Name string
	Fn   func(map[string]string) (map[string]string, error)
}

// Chain executes a sequence of steps over a secrets map, passing the output
// of each step as the input to the next.
type Chain struct {
	steps []ChainStep
}

// NewChain creates a new Chain with no steps.
func NewChain() *Chain {
	return &Chain{}
}

// Add appends a named step to the chain.
func (c *Chain) Add(name string, fn func(map[string]string) (map[string]string, error)) *Chain {
	c.steps = append(c.steps, ChainStep{Name: name, Fn: fn})
	return c
}

// Run executes all steps in order. If any step returns an error the chain
// halts and returns the step name alongside the error.
func (c *Chain) Run(secrets map[string]string) (map[string]string, error) {
	current := copyMap(secrets)
	for _, step := range c.steps {
		var err error
		current, err = step.Fn(current)
		if err != nil {
			return nil, fmt.Errorf("chain step %q: %w", step.Name, err)
		}
	}
	return current, nil
}

// Steps returns the names of all registered steps.
func (c *Chain) Steps() []string {
	names := make([]string, len(c.steps))
	for i, s := range c.steps {
		names[i] = s.Name
	}
	return names
}

func copyMap(m map[string]string) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}
