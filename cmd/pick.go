package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/envchain/envchain/internal/pick"
	"github.com/spf13/cobra"
)

func init() {
	var keys []string
	var pattern string

	cmd := &cobra.Command{
		Use:   "pick",
		Short: "Select a subset of keys from an env map",
		Long:  "Pick reads a JSON env map from stdin and outputs only the selected keys.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runPick(cmd, keys, pattern)
		},
	}

	cmd.Flags().StringSliceVarP(&keys, "key", "k", nil, "explicit keys to include (comma-separated)")
	cmd.Flags().StringVarP(&pattern, "pattern", "p", "", "regex pattern to match keys")

	rootCmd.AddCommand(cmd)
}

func runPick(cmd *cobra.Command, keys []string, pattern string) error {
	var src map[string]string
	if err := json.NewDecoder(os.Stdin).Decode(&src); err != nil {
		return fmt.Errorf("pick: failed to decode stdin: %w", err)
	}

	// filter out empty strings that may come from cobra flag parsing
	filtered := keys[:0]
	for _, k := range keys {
		if strings.TrimSpace(k) != "" {
			filtered = append(filtered, k)
		}
	}

	out, err := pick.Pick(src, pick.Options{
		Keys:    filtered,
		Pattern: pattern,
	})
	if err != nil {
		return err
	}

	return json.NewEncoder(cmd.OutOrStdout()).Encode(out)
}
