package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"

	"github.com/spf13/cobra"

	"envchain/internal/redact"
)

var (
	redactPatterns  []string
	redactReplace   string
	redactInputFile string
)

func init() {
	redactCmd := &cobra.Command{
		Use:   "redact",
		Short: "Redact sensitive values from an env map",
		RunE:  runRedact,
	}
	redactCmd.Flags().StringArrayVarP(&redactPatterns, "pattern", "p", []string{`password|secret|token`}, "regex patterns to match sensitive keys (can repeat)")
	redactCmd.Flags().StringVar(&redactReplace, "replace", "", "replacement string (default: [REDACTED])")
	redactCmd.Flags().StringVarP(&redactInputFile, "file", "f", "", "JSON file containing env map (reads stdin if omitted)")
	rootCmd.AddCommand(redactCmd)
}

func runRedact(cmd *cobra.Command, args []string) error {
	var raw []byte
	var err error

	if redactInputFile != "" {
		raw, err = os.ReadFile(redactInputFile)
	} else {
		raw, err = readStdin()
	}
	if err != nil {
		return fmt.Errorf("redact: read input: %w", err)
	}

	var envMap map[string]string
	if err := json.Unmarshal(raw, &envMap); err != nil {
		return fmt.Errorf("redact: parse JSON: %w", err)
	}

	var rules []redact.Rule
	for i, p := range redactPatterns {
		re, err := regexp.Compile(p)
		if err != nil {
			return fmt.Errorf("redact: invalid pattern %q: %w", p, err)
		}
		rules = append(rules, redact.Rule{
			Name:    fmt.Sprintf("rule-%d", i),
			Pattern: re,
			Replace: redactReplace,
		})
	}

	r, err := redact.New(redact.Options{Rules: rules})
	if err != nil {
		return err
	}

	out := r.RedactMap(envMap)
	enc := json.NewEncoder(cmd.OutOrStdout())
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}

func readStdin() ([]byte, error) {
	var buf []byte
	tmp := make([]byte, 512)
	for {
		n, err := os.Stdin.Read(tmp)
		buf = append(buf, tmp[:n]...)
		if err != nil {
			break
		}
	}
	return buf, nil
}
