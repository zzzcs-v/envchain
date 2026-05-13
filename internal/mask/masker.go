package mask

import (
	"fmt"
	"regexp"
	"strings"
)

// Strategy defines how values are masked.
type Strategy int

const (
	StrategyRedact  Strategy = iota // replaces entire value with ***
	StrategyPartial                 // shows first/last N chars
	StrategyHash                    // shows a short hash
)

// Options configures masking behaviour.
type Options struct {
	Strategy     Strategy
	VisibleChars int    // used by StrategyPartial
	KeyPattern   string // regex; only mask matching keys
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		Strategy:     StrategyRedact,
		VisibleChars: 3,
		KeyPattern:   `(?i)(secret|password|token|key|pass|pwd|auth|credential)`,
	}
}

// Masker masks sensitive environment variable values.
type Masker struct {
	opts  Options
	keyRe *regexp.Regexp
}

// New creates a Masker with the given options.
func New(opts Options) (*Masker, error) {
	var re *regexp.Regexp
	if opts.KeyPattern != "" {
		var err error
		re, err = regexp.Compile(opts.KeyPattern)
		if err != nil {
			return nil, err
		}
	}
	return &Masker{opts: opts, keyRe: re}, nil
}

// MaskMap returns a copy of vars with sensitive values masked.
func (m *Masker) MaskMap(vars map[string]string) map[string]string {
	out := make(map[string]string, len(vars))
	for k, v := range vars {
		if m.isSensitive(k) {
			out[k] = m.maskValue(v)
		} else {
			out[k] = v
		}
	}
	return out
}

// MaskValue masks a single value unconditionally.
func (m *Masker) MaskValue(v string) string {
	return m.maskValue(v)
}

func (m *Masker) isSensitive(key string) bool {
	if m.keyRe == nil {
		return true
	}
	return m.keyRe.MatchString(key)
}

func (m *Masker) maskValue(v string) string {
	switch m.opts.Strategy {
	case StrategyPartial:
		return partialMask(v, m.opts.VisibleChars)
	case StrategyHash:
		return hashMask(v)
	default:
		return "***"
	}
}

func partialMask(v string, n int) string {
	if len(v) <= n*2 {
		return strings.Repeat("*", len(v))
	}
	return v[:n] + strings.Repeat("*", len(v)-n*2) + v[len(v)-n:]
}

func hashMask(v string) string {
	if len(v) == 0 {
		return "***"
	}
	var h uint32
	for _, c := range v {
		h = h*31 + uint32(c)
	}
	return fmt.Sprintf("#%06x", h&0xFFFFFF)
}
