package audit

import "fmt"

// ActionExport is the audit action for export operations.
const (
	ActionExport   = "export"
	ActionValidate = "validate"
	ActionDiff     = "diff"
)

// Middleware wraps a Logger and provides helpers for common envchain actions.
type Middleware struct {
	logger *Logger
	enabled bool
}

// NewMiddleware creates a Middleware. If path is empty, logging is a no-op.
func NewMiddleware(path string) (*Middleware, error) {
	if path == "" {
		return &Middleware{enabled: false}, nil
	}
	l, err := NewLogger(path)
	if err != nil {
		return nil, fmt.Errorf("audit middleware: %w", err)
	}
	return &Middleware{logger: l, enabled: true}, nil
}

// LogExport records an export action.
func (m *Middleware) LogExport(context, format string, vars []string) error {
	if !m.enabled {
		return nil
	}
	return m.logger.Log(ActionExport, context, vars, map[string]string{"format": format})
}

// LogValidate records a validate action.
func (m *Middleware) LogValidate(context string, vars []string) error {
	if !m.enabled {
		return nil
	}
	return m.logger.Log(ActionValidate, context, vars, nil)
}

// LogDiff records a diff action.
func (m *Middleware) LogDiff(contextA, contextB string) error {
	if !m.enabled {
		return nil
	}
	return m.logger.Log(ActionDiff, contextA, nil, map[string]string{"compare_to": contextB})
}

// Close closes the underlying logger if enabled.
func (m *Middleware) Close() error {
	if !m.enabled || m.logger == nil {
		return nil
	}
	return m.logger.Close()
}
