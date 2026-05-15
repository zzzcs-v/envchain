package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/user/envchain/internal/defaults"
)

func init() {
	var override bool

	cmd := &cobra.Command{
		Use:   "defaults",
		Short: "Apply default values to a JSON env map from stdin",
		Long: `Reads a JSON object from stdin and applies default key=value pairs.
Existing non-empty values are preserved unless --override is set.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDefaults(cmd, args, override)
		},
	}

	cmd.Flags().BoolVar(&override, "override", false, "Replace existing non-empty values")
	cmd.Flags().StringArrayP("set", "s", nil, "Default entry as KEY=VALUE (repeatable)")
	_ = cmd.MarkFlagRequired("set")

	rootCmd.AddCommand(cmd)
}

func runDefaults(cmd *cobra.Command, _ []string, override bool) error {
	pairs, _ := cmd.Flags().GetStringArray("set")

	var entries []defaults.Entry
	for _, p := range pairs {
		parts := strings.SplitN(p, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid pair %q: expected KEY=VALUE", p)
		}
		entries = append(entries, defaults.Entry{
			Key:      parts[0],
			Value:    parts[1],
			Override: override,
		})
	}

	var src map[string]string
	if err := json.NewDecoder(os.Stdin).Decode(&src); err != nil {
		return fmt.Errorf("failed to decode stdin: %w", err)
	}

	res, err := defaults.Apply(src, entries)
	if err != nil {
		return err
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(src); err != nil {
		return fmt.Errorf("failed to encode output: %w", err)
	}

	fmt.Fprintf(os.Stderr, "# %s\n", res.Summary())
	return nil
}
