package chain

import "fmt"

// Chain represents an ordered sequence of named contexts to apply in order.
type Chain struct {
	Name     string
	Contexts []string
}

// Store holds named chains.
type Store struct {
	chains map[string]Chain
}

// New returns an empty Store.
func New() *Store {
	return &Store{chains: make(map[string]Chain)}
}

// Set registers a named chain with the given ordered context names.
func (s *Store) Set(name string, contexts []string) error {
	if name == "" {
		return fmt.Errorf("chain name must not be empty")
	}
	if len(contexts) == 0 {
		return fmt.Errorf("chain %q must have at least one context", name)
	}
	seen := make(map[string]bool)
	for _, c := range contexts {
		if c == "" {
			return fmt.Errorf("chain %q contains an empty context name", name)
		}
		if seen[c] {
			return fmt.Errorf("chain %q contains duplicate context %q", name, c)
		}
		seen[c] = true
	}
	s.chains[name] = Chain{Name: name, Contexts: contexts}
	return nil
}

// Get retrieves a chain by name.
func (s *Store) Get(name string) (Chain, error) {
	c, ok := s.chains[name]
	if !ok {
		return Chain{}, fmt.Errorf("chain %q not found", name)
	}
	return c, nil
}

// Delete removes a chain by name.
func (s *Store) Delete(name string) error {
	if _, ok := s.chains[name]; !ok {
		return fmt.Errorf("chain %q not found", name)
	}
	delete(s.chains, name)
	return nil
}

// List returns all chains sorted by name.
func (s *Store) List() []Chain {
	out := make([]Chain, 0, len(s.chains))
	for _, c := range s.chains {
		out = append(out, c)
	}
	sortChains(out)
	return out
}

func sortChains(chains []Chain) {
	for i := 1; i < len(chains); i++ {
		for j := i; j > 0 && chains[j].Name < chains[j-1].Name; j-- {
			chains[j], chains[j-1] = chains[j-1], chains[j]
		}
	}
}
