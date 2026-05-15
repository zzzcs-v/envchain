package transform

import (
	"fmt"
	"strings"
)

// Op represents a transformation operation.
type Op string

const (
	OpUppercase  Op = "uppercase"
	OpLowercase  Op = "lowercase"
	OpTrimSpace  Op = "trimspace"
	OpBase64Enc  Op = "base64enc"
	OpBase64Dec  Op = "base64dec"
)

// Options controls how transformations are applied.
type Options struct {
	// Keys restricts transformation to specific keys; empty means all.
	Keys []string
	// Op is the transformation to apply.
	Op Op
}

// Result holds the transformed map and a summary of changes.
type Result struct {
	Vars    map[string]string
	Changed int
}

// Apply applies the configured transformation to src and returns a new map.
func Apply(src map[string]string, opts Options) (Result, error) {
	if src == nil {
		return Result{Vars: map[string]string{}}, nil
	}

	fn, err := opFunc(opts.Op)
	if err != nil {
		return Result{}, err
	}

	allow := toSet(opts.Keys)
	out := make(map[string]string, len(src))
	changed := 0

	for k, v := range src {
		if len(allow) > 0 && !allow[k] {
			out[k] = v
			continue
		}
		nv := fn(v)
		if nv != v {
			changed++
		}
		out[k] = nv
	}

	return Result{Vars: out, Changed: changed}, nil
}

func opFunc(op Op) (func(string) string, error) {
	switch op {
	case OpUppercase:
		return strings.ToUpper, nil
	case OpLowercase:
		return strings.ToLower, nil
	case OpTrimSpace:
		return strings.TrimSpace, nil
	case OpBase64Enc:
		return b64enc, nil
	case OpBase64Dec:
		return b64dec, nil
	default:
		return nil, fmt.Errorf("unknown transform op %q", op)
	}
}

func toSet(keys []string) map[string]bool {
	s := make(map[string]bool, len(keys))
	for _, k := range keys {
		s[k] = true
	}
	return s
}
