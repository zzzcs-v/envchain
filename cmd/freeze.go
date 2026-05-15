package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/envchain/envchain/internal/freeze"
	"github.com/spf13/cobra"
)

func init() {
	var pairsFlag []string
	var diffFlag []string

	cmd := &cobra.Command{
		Use:   "freeze",
		Short: "Freeze an env map and optionally diff against a live set",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runFreeze(pairsFlag, diffFlag)
		},
	}

	cmd.Flags().StringArrayVar(&pairsFlag, "set", nil, "KEY=VALUE pairs to freeze (required)")
	cmd.Flags().StringArrayVar(&diffFlag, "diff", nil, "KEY=VALUE pairs to diff against the frozen snapshot")
	_ = cmd.MarkFlagRequired("set")

	rootCmd.AddCommand(cmd)
}

func runFreeze(pairs []string, diffPairs []string) error {
	source, err := parseFreezeKV(pairs)
	if err != nil {
		return fmt.Errorf("--set: %w", err)
	}

	f, err := freeze.New(source)
	if err != nil {
		return err
	}

	if len(diffPairs) == 0 {
		out, _ := json.MarshalIndent(f.Snapshot(), "", "  ")
		fmt.Fprintln(os.Stdout, string(out))
		return nil
	}

	live, err := parseFreezeKV(diffPairs)
	if err != nil {
		return fmt.Errorf("--diff: %w", err)
	}

	changed := f.DiffFrom(live)
	if len(changed) == 0 {
		fmt.Fprintln(os.Stdout, "no drift detected")
		return nil
	}

	fmt.Fprintln(os.Stdout, "drifted keys:")
	for _, k := range changed {
		fmt.Fprintf(os.Stdout, "  %s\n", k)
	}
	return nil
}

func parseFreezeKV(pairs []string) (map[string]string, error) {
	out := make(map[string]string, len(pairs))
	for _, p := range pairs {
		parts := strings.SplitN(p, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid pair %q: expected KEY=VALUE", p)
		}
		out[parts[0]] = parts[1]
	}
	return out, nil
}
