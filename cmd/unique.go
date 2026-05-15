package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envchain/internal/unique"
)

func init() {
	var (
		keys          []string
		caseSensitive bool
		showDropped   bool
	)

	cmd := &cobra.Command{
		Use:   "unique",
		Short: "Remove entries with duplicate values from a JSON env map (stdin)",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runUnique(os.Stdin, keys, caseSensitive, showDropped)
		},
	}

	cmd.Flags().StringSliceVarP(&keys, "keys", "k", nil, "restrict uniqueness check to these keys")
	cmd.Flags().BoolVar(&caseSensitive, "case-sensitive", false, "compare values case-sensitively")
	cmd.Flags().BoolVar(&showDropped, "show-dropped", false, "print dropped keys instead of kept keys")

	rootCmd.AddCommand(cmd)
}

func runUnique(r interface{ Read([]byte) (int, error) }, keys []string, caseSensitive, showDropped bool) error {
	var src map[string]string
	if err := json.NewDecoder(r.(interface {
		Read([]byte) (int, error)
		Close() error
	})).Decode(&src); err != nil {
		if f, ok := r.(*os.File); ok {
			_ = f
		}
		return fmt.Errorf("decode input: %w", err)
	}

	res, err := unique.Unique(src, unique.Options{
		Keys:          keys,
		CaseSensitive: caseSensitive,
	})
	if err != nil {
		return err
	}

	out := res.Kept
	if showDropped {
		out = res.Dropped
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}
