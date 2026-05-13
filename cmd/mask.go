package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"

	"github.com/spf13/cobra"

	"envchain/internal/mask"
)

var (
	maskStrategy string
	maskVisible  int
	maskPattern  string
)

func init() {
	maskCmd := &cobra.Command{
		Use:   "mask",
		Short: "Print environment variables with sensitive values masked",
		RunE:  runMask,
	}
	maskCmd.Flags().StringVar(&maskStrategy, "strategy", "redact", "masking strategy: redact, partial, hash")
	maskCmd.Flags().IntVar(&maskVisible, "visible", 3, "visible chars on each side (partial strategy)")
	maskCmd.Flags().StringVar(&maskPattern, "pattern", mask.DefaultOptions().KeyPattern, "regex pattern for sensitive key names")
	rootCmd.AddCommand(maskCmd)
}

func runMask(cmd *cobra.Command, args []string) error {
	opts := mask.Options{
		VisibleChars: maskVisible,
		KeyPattern:   maskPattern,
	}
	switch maskStrategy {
	case "partial":
		opts.Strategy = mask.StrategyPartial
	case "hash":
		opts.Strategy = mask.StrategyHash
	default:
		opts.Strategy = mask.StrategyRedact
	}

	m, err := mask.New(opts)
	if err != nil {
		return fmt.Errorf("invalid pattern: %w", err)
	}

	envVars := map[string]string{}
	for _, e := range os.Environ() {
		for i := 0; i < len(e); i++ {
			if e[i] == '=' {
				envVars[e[:i]] = e[i+1:]
				break
			}
		}
	}

	masked := m.MaskMap(envVars)

	keys := make([]string, 0, len(masked))
	for k := range masked {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	enc := json.NewEncoder(cmd.OutOrStdout())
	enc.SetIndent("", "  ")
	ordered := make(map[string]string, len(keys))
	for _, k := range keys {
		ordered[k] = masked[k]
	}
	return enc.Encode(ordered)
}
