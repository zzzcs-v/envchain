package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// EnvConfig holds environment variables for a named context
type EnvConfig struct {
	Name    string            `json:"name"`
	Context string            `json:"context"`
	Vars    map[string]string `json:"vars"`
}

// ChainConfig represents the top-level envchain config file
type ChainConfig struct {
	Version  string      `json:"version"`
	Contexts []EnvConfig `json:"contexts"`
}

const defaultConfigFile = ".envchain.json"

// Load reads and parses the envchain config file from the given path.
// If path is empty, it looks for .envchain.json in the current directory.
func Load(path string) (*ChainConfig, error) {
	if path == "" {
		path = defaultConfigFile
	}

	abs, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("resolving config path: %w", err)
	}

	data, err := os.ReadFile(abs)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("config file not found: %s", abs)
		}
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	var cfg ChainConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config file: %w", err)
	}

	if err := validate(&cfg); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &cfg, nil
}

// GetContext returns the EnvConfig for the given context name.
func (c *ChainConfig) GetContext(name string) (*EnvConfig, error) {
	for i := range c.Contexts {
		if c.Contexts[i].Context == name {
			return &c.Contexts[i], nil
		}
	}
	return nil, fmt.Errorf("context %q not found in config", name)
}

func validate(cfg *ChainConfig) error {
	if cfg.Version == "" {
		return fmt.Errorf("missing required field: version")
	}
	seen := make(map[string]bool)
	for _, ctx := range cfg.Contexts {
		if ctx.Context == "" {
			return fmt.Errorf("context entry missing required field: context")
		}
		if seen[ctx.Context] {
			return fmt.Errorf("duplicate context name: %s", ctx.Context)
		}
		seen[ctx.Context] = true
	}
	return nil
}
