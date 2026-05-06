package context

import (
	"fmt"
	"os"
	"strings"
)

// Context holds the resolved environment variables for a given context name.
type Context struct {
	Name string
	Vars map[string]string
}

// Resolver resolves environment variables for a named context,
// supporting inheritance via a "extends" chain.
type Resolver struct {
	contexts map[string]rawContext
}

type rawContext struct {
	Extends string
	Vars    map[string]string
}

// NewResolver creates a Resolver from a map of context definitions.
func NewResolver(defs map[string]map[string]string, extends map[string]string) *Resolver {
	r := &Resolver{contexts: make(map[string]rawContext)}
	for name, vars := range defs {
		r.contexts[name] = rawContext{
			Extends: extends[name],
			Vars:    vars,
		}
	}
	return r
}

// Resolve returns a Context with all variables merged, respecting the extends chain.
func (r *Resolver) Resolve(name string) (*Context, error) {
	seen := map[string]bool{}
	vars, err := r.merge(name, seen)
	if err != nil {
		return nil, err
	}
	return &Context{Name: name, Vars: vars}, nil
}

func (r *Resolver) merge(name string, seen map[string]bool) (map[string]string, error) {
	if seen[name] {
		return nil, fmt.Errorf("circular extends detected at context %q", name)
	}
	seen[name] = true

	ctx, ok := r.contexts[name]
	if !ok {
		return nil, fmt.Errorf("context %q not found", name)
	}

	base := map[string]string{}
	if ctx.Extends != "" {
		var err error
		base, err = r.merge(ctx.Extends, seen)
		if err != nil {
			return nil, err
		}
	}

	for k, v := range ctx.Vars {
		base[k] = expandEnv(v)
	}
	return base, nil
}

// expandEnv replaces ${VAR} and $VAR references with values from the OS environment.
func expandEnv(val string) string {
	return os.Expand(val, func(key string) string {
		if v, ok := os.LookupEnv(key); ok {
			return v
		}
		return "$" + key
	})
}

// ToEnvSlice returns the context variables as a slice of KEY=VALUE strings.
func (c *Context) ToEnvSlice() []string {
	out := make([]string, 0, len(c.Vars))
	for k, v := range c.Vars {
		out = append(out, strings.ToUpper(k)+"="+v)
	}
	return out
}
