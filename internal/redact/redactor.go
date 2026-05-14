package redact

import (
	"fmt"
	"regexp"
	"strings"
)

// Rule defines a single redaction rule.
type Rule struct {
	Name    string
	Pattern *regexp.Regexp
	Replace string
}

// Options configures the Redactor.
type Options struct {
	Rules      []Rule
	Placeholder string // default: "[REDACTED]"
}

// Redactor applies redaction rules to maps and strings.
type Redactor struct {
	opts Options
}

// New creates a Redactor with the given options.
func New(opts Options) (*Redactor, error) {
	if opts.Placeholder == "" {
		opts.Placeholder = "[REDACTED]"
	}
	for _, r := range opts.Rules {
		if r.Pattern == nil {
			return nil, fmt.Errorf("redact: rule %q has nil pattern", r.Name)
		}
	}
	return &Redactor{opts: opts}, nil
}

// RedactMap returns a copy of m with values redacted according to rules.
// Keys matching any rule pattern have their value replaced with the placeholder
// (or rule-specific replacement if set).
func (r *Redactor) RedactMap(m map[string]string) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = r.redactValue(k, v)
	}
	return out
}

// RedactString replaces all matches of any rule pattern within s.
func (r *Redactor) RedactString(s string) string {
	for _, rule := range r.opts.Rules {
		repl := rule.Replace
		if repl == "" {
			repl = r.opts.Placeholder
		}
		s = rule.Pattern.ReplaceAllString(s, repl)
	}
	return s
}

func (r *Redactor) redactValue(key, value string) string {
	for _, rule := range r.opts.Rules {
		if rule.Pattern.MatchString(strings.ToLower(key)) {
			repl := rule.Replace
			if repl == "" {
				repl = r.opts.Placeholder
			}
			return repl
		}
	}
	return value
}
