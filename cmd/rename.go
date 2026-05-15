package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envchain/internal/rename"
)

func init() {
	var (
		pattern         string
		template        string
		skipConflicts   bool
		errorOnConflict bool
		showResults     bool
	)

	cmd := &cobra.Command{
		Use:   "rename",
		Short: "Rename env keys using a regex pattern",
		Example: `  envchain rename --from '^DB_' --to 'PG_' --vars 'DB_HOST=localhost,DB_PORT=5432'`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRename(cmd, pattern, template, skipConflicts, errorOnConflict, showResults)
		},
	}

	cmd.Flags().StringVar(&pattern, "from", "", "regex pattern to match keys (required)")
	cmd.Flags().StringVar(&template, "to", "", "replacement template for matched keys (required)")
	cmd.Flags().BoolVar(&skipConflicts, "skip-conflicts", false, "silently skip renames that would conflict")
	cmd.Flags().BoolVar(&errorOnConflict, "error-on-conflict", false, "return an error on key conflicts")
	cmd.Flags().BoolVar(&showResults, "show-results", false, "print rename results to stderr")
	_ = cmd.MarkFlagRequired("from")
	_ = cmd.MarkFlagRequired("to")

	rootCmd.AddCommand(cmd)
}

func runRename(cmd *cobra.Command, pattern, template string, skipConflicts, errorOnConflict, showResults bool) error {
	// Read vars from stdin as JSON for pipeline use.
	var src map[string]string
	if err := json.NewDecoder(os.Stdin).Decode(&src); err != nil {
		return fmt.Errorf("rename: failed to decode stdin JSON: %w", err)
	}

	out, results, err := rename.Rename(src, rename.Options{
		FromPattern:     pattern,
		ToTemplate:      template,
		SkipConflicts:   skipConflicts,
		ErrorOnConflict: errorOnConflict,
	})
	if err != nil {
		return err
	}

	if showResults {
		for _, r := range results {
			if r.Skipped {
				fmt.Fprintf(os.Stderr, "skipped: %s -> %s (conflict)\n", r.OldKey, r.NewKey)
			} else {
				fmt.Fprintf(os.Stderr, "renamed: %s -> %s\n", r.OldKey, r.NewKey)
			}
		}
	}

	return json.NewEncoder(cmd.OutOrStdout()).Encode(out)
}
