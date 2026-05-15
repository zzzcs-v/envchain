package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/envchain/envchain/internal/truncate"
)

func init() {
	var maxLen int
	var suffix string
	var keys []string

	cmd := &cobra.Command{
		Use:   "truncate",
		Short: "Truncate long env var values to a maximum length",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runTruncate(cmd, maxLen, suffix, keys)
		},
	}

	cmd.Flags().IntVar(&maxLen, "max-len", 64, "maximum value length in runes (0 = unlimited)")
	cmd.Flags().StringVar(&suffix, "suffix", "", "suffix to append when a value is truncated (e.g. \"...\")")
	cmd.Flags().StringSliceVar(&keys, "keys", nil, "restrict truncation to these keys (default: all)")

	rootCmd.AddCommand(cmd)
}

func runTruncate(cmd *cobra.Command, maxLen int, suffix string, keys []string) error {
	var src map[string]string
	dec := json.NewDecoder(os.Stdin)
	if err := dec.Decode(&src); err != nil {
		return fmt.Errorf("truncate: failed to decode stdin: %w", err)
	}

	out, results, err := truncate.Truncate(src, truncate.Options{
		MaxLen: maxLen,
		Suffix: suffix,
		Keys:   keys,
	})
	if err != nil {
		return err
	}

	if len(results) > 0 {
		for _, r := range results {
			fmt.Fprintf(cmd.ErrOrStderr(), "truncated %s (%d → %d runes)\n",
				r.Key, len([]rune(r.Original)), len([]rune(r.Truncated)))
		}
	}

	enc := json.NewEncoder(cmd.OutOrStdout())
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}
