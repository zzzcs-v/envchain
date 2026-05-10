package prompt

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"syscall"

	"golang.org/x/term"
)

// Prompter handles interactive user input for CLI prompts.
type Prompter struct {
	in  io.Reader
	out io.Writer
}

// New returns a Prompter using stdin/stdout by default.
func New() *Prompter {
	return &Prompter{in: os.Stdin, out: os.Stdout}
}

// NewWithIO returns a Prompter with custom reader/writer (useful for testing).
func NewWithIO(in io.Reader, out io.Writer) *Prompter {
	return &Prompter{in: in, out: out}
}

// Ask prompts the user for a plain-text value.
func (p *Prompter) Ask(label string) (string, error) {
	fmt.Fprintf(p.out, "%s: ", label)
	scanner := bufio.NewScanner(p.in)
	if scanner.Scan() {
		return strings.TrimSpace(scanner.Text()), nil
	}
	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("prompt read error: %w", err)
	}
	return "", fmt.Errorf("no input received for %q", label)
}

// AskSecret prompts the user for a secret value (no echo) using the terminal.
// Falls back to plain Ask when not connected to a real terminal (e.g. in tests).
func (p *Prompter) AskSecret(label string) (string, error) {
	fmt.Fprintf(p.out, "%s (hidden): ", label)
	if p.in == os.Stdin && term.IsTerminal(int(syscall.Stdin)) {
		bytes, err := term.ReadPassword(int(syscall.Stdin))
		fmt.Fprintln(p.out)
		if err != nil {
			return "", fmt.Errorf("failed to read secret: %w", err)
		}
		return strings.TrimSpace(string(bytes)), nil
	}
	// Fallback for non-terminal environments.
	scanner := bufio.NewScanner(p.in)
	if scanner.Scan() {
		fmt.Fprintln(p.out)
		return strings.TrimSpace(scanner.Text()), nil
	}
	return "", fmt.Errorf("no secret input received for %q", label)
}

// Confirm asks a yes/no question and returns true if the user answers 'y' or 'yes'.
func (p *Prompter) Confirm(label string) (bool, error) {
	fmt.Fprintf(p.out, "%s [y/N]: ", label)
	scanner := bufio.NewScanner(p.in)
	if scanner.Scan() {
		answer := strings.ToLower(strings.TrimSpace(scanner.Text()))
		return answer == "y" || answer == "yes", nil
	}
	if err := scanner.Err(); err != nil {
		return false, fmt.Errorf("confirm read error: %w", err)
	}
	return false, nil
}
