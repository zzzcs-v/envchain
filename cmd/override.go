package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envchain/internal/override"
)

func init() {
	var (
		pairs      []string
		allowEmpty bool
		strict     bool
	)

	cmd := &cobra.Command{
		Use:   "override",
		Short: "Apply key=value overrides to an env map and print the result",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runOverride(pairs, allowEmpty, strict)
		},
	}

	cmd.Flags().StringArrayVarP(&pairs, "set", "s", nil, "KEY=VALUE pairs to override (repeatable)")
	cmd.Flags().BoolVar(&allowEmpty, "allow-empty", false, "permit overriding keys with empty values")
	cmd.Flags().BoolVar(&strict, "strict", false, "fail if an override key does not exist in the base map")
	_ = cmd.MarkFlagRequired("set")

	rootCmd.AddCommand(cmd)
}

func runOverride(pairs []string, allowEmpty, strict bool) error {
	if len(pairs) == 0 {
		return fmt.Errorf("override: at least one --set pair is required")
	}

	// Start with env vars from OS as the base map.
	base := map[string]string{}
	for _, e := range os.Environ() {
		entry, err := override.ParsePair(e)
		if err == nil {
			base[entry.Key] = entry.Value
		}
	}

	entries := make([]override.Entry, 0, len(pairs))
	for _, p := range pairs {
		e, err := override.ParsePair(p)
		if err != nil {
			return err
		}
		entries = append(entries, e)
	}

	opts := override.Options{AllowEmpty: allowEmpty, Strict: strict}
	if err := override.Apply(base, entries, opts); err != nil {
		return err
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(base)
}
