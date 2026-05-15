package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envchain/internal/omit"
)

var (
	omitKeys    []string
	omitPattern string
	omitPrefix  string
)

func init() {
	cmd := &cobra.Command{
		Use:   "omit",
		Short: "Remove keys from an env map by name, pattern, or prefix",
		RunE:  runOmit,
	}
	cmd.Flags().StringSliceVarP(&omitKeys, "key", "k", nil, "explicit keys to remove (comma-separated)")
	cmd.Flags().StringVarP(&omitPattern, "pattern", "p", "", "regex pattern; matching keys are removed")
	cmd.Flags().StringVar(&omitPrefix, "prefix", "", "remove all keys that start with this prefix")
	rootCmd.AddCommand(cmd)
}

func runOmit(cmd *cobra.Command, _ []string) error {
	if omitPattern == "" && len(omitKeys) == 0 && omitPrefix == "" {
		return fmt.Errorf("omit: at least one of --key, --pattern, or --prefix is required")
	}

	var src map[string]string
	if err := json.NewDecoder(os.Stdin).Decode(&src); err != nil {
		return fmt.Errorf("omit: failed to decode stdin: %w", err)
	}

	out, err := omit.Omit(src, omit.Options{
		Keys:    omitKeys,
		Pattern: omitPattern,
		Prefix:  omitPrefix,
	})
	if err != nil {
		return err
	}

	return json.NewEncoder(cmd.OutOrStdout()).Encode(out)
}
