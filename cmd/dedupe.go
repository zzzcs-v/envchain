package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/envchain/envchain/internal/dedupe"
	"github.com/spf13/cobra"
)

func init() {
	var strategy string
	var restrictKeys []string

	cmd := &cobra.Command{
		Use:   "dedupe",
		Short: "Remove duplicate keys across environment variable sources",
		Long: `Reads one or more JSON env maps from --source flags and removes duplicate
keys, keeping either the first or last occurrence depending on --strategy.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			sources, err := cmd.Flags().GetStringArray("source")
			if err != nil {
				return err
			}
			var maps []map[string]string
			for _, s := range sources {
				var m map[string]string
				if err := json.Unmarshal([]byte(s), &m); err != nil {
					return fmt.Errorf("invalid source JSON %q: %w", s, err)
				}
				maps = append(maps, m)
			}

			var strat dedupe.Strategy
			switch strategy {
			case "first":
				strat = dedupe.KeepFirst
			case "last":
				strat = dedupe.KeepLast
			default:
				return fmt.Errorf("unknown strategy %q: must be 'first' or 'last'", strategy)
			}

			res, err := dedupe.Dedupe(maps, dedupe.Options{
				Strategy:       strat,
				RestrictToKeys: restrictKeys,
			})
			if err != nil {
				return err
			}

			out, err := json.MarshalIndent(res.Vars, "", "  ")
			if err != nil {
				return err
			}
			fmt.Fprintln(os.Stdout, string(out))
			if len(res.Removed) > 0 {
				fmt.Fprintf(os.Stderr, "info: %s\n", res.Summary())
			}
			return nil
		},
	}

	cmd.Flags().StringArray("source", nil, "JSON env map (repeatable)")
	cmd.Flags().StringVar(&strategy, "strategy", "first", "dedup strategy: first|last")
	cmd.Flags().StringSliceVar(&restrictKeys, "keys", nil, "restrict dedup to these keys only")
	_ = cmd.MarkFlagRequired("source")

	rootCmd.AddCommand(cmd)
}
