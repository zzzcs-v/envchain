package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"envchain/internal/prefix"
)

func init() {
	var (
		pfx      string
		strip    bool
		overwrite bool
	)

	cmd := &cobra.Command{
		Use:   "prefix",
		Short: "Add or strip a prefix from env var keys",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runPrefix(cmd, pfx, strip, overwrite)
		},
	}

	cmd.Flags().StringVarP(&pfx, "prefix", "p", "", "prefix string to add or strip (required)")
	cmd.Flags().BoolVar(&strip, "strip", false, "strip prefix instead of adding it")
	cmd.Flags().BoolVar(&overwrite, "overwrite", false, "re-prefix keys that already have the prefix")
	_ = cmd.MarkFlagRequired("prefix")

	rootCmd.AddCommand(cmd)
}

func runPrefix(cmd *cobra.Command, pfx string, strip, overwrite bool) error {
	src := map[string]string{}
	for _, e := range os.Environ() {
		parts := strings.SplitN(e, "=", 2)
		if len(parts) == 2 {
			src[parts[0]] = parts[1]
		}
	}

	out, res, err := prefix.Apply(src, prefix.Options{
		Prefix:    pfx,
		Strip:     strip,
		Overwrite: overwrite,
	})
	if err != nil {
		return fmt.Errorf("prefix: %w", err)
	}

	enc := json.NewEncoder(cmd.OutOrStdout())
	enc.SetIndent("", "  ")
	if err := enc.Encode(out); err != nil {
		return fmt.Errorf("prefix: encode: %w", err)
	}

	fmt.Fprintln(cmd.ErrOrStderr(), res.Summary())
	return nil
}
