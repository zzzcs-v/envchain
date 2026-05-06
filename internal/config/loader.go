package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config is the top-level structure for an envchain config file.
type Config struct {
	Version  string             `yaml:"version"`
	Contexts map[string]CtxDef `yaml:"contexts"`
}

// CtxDef defines a single context entry in the config file.
type CtxDef struct {
	Extends string            `yaml:"extends,omitempty"`
	Vars    map[string]string `yaml:"vars"`
}

// Load reads and parses an envchain YAML config file from the given path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}

	if err := validate(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// ToResolverInputs extracts the defs and extends maps needed by context.NewResolver.
func (c *Config) ToResolverInputs() (defs map[string]map[string]string, extends map[string]string) {
	defs = make(map[string]map[string]string, len(c.Contexts))
	extends = make(map[string]string, len(c.Contexts))
	for name, ctx := range c.Contexts {
		defs[name] = ctx.Vars
		if ctx.Extends != "" {
			extends[name] = ctx.Extends
		}
	}
	return
}

func validate(cfg *Config) error {
	if cfg.Version == "" {
		return fmt.Errorf("config missing required field: version")
	}
	names := make(map[string]bool, len(cfg.Contexts))
	for name := range cfg.Contexts {
		if names[name] {
			return fmt.Errorf("duplicate context name: %q", name)
		}
		names[name] = true
	}
	for name, ctx := range cfg.Contexts {
		if ctx.Extends != "" {
			if _, ok := cfg.Contexts[ctx.Extends]; !ok {
				return fmt.Errorf("context %q extends unknown context %q", name, ctx.Extends)
			}
		}
	}
	return nil
}
