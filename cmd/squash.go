package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envchain/internal/squash"
)

var squashCmd = &cobra.Command{
	Use:   "squash",
	Short: "Merge multiple env maps into one, resolving conflicts by strategy",
	RunE:  runSquash,
}

var squashStrategy string
var squashSources []string

func init() {
	squashCmd.Flags().StringVarP(&squashStrategy, "strategy", "s", "keep-last",
		"conflict resolution strategy: keep-first | keep-last | error")
	squashCmd.Flags().StringArrayVarP(&squashSources, "source", "S", nil,
		"JSON env map (may be repeated); reads stdin if omitted")
	rootCmd.AddCommand(squashCmd)
}

func runSquash(cmd *cobra.Command, _ []string) error {
	strat, err := squash.ParseStrategy(squashStrategy)
	if err != nil {
		return err
	}

	var sources []map[string]string

	if len(squashSources) == 0 {
		var m map[string]string
		if err := json.NewDecoder(os.Stdin).Decode(&m); err != nil {
			return fmt.Errorf("squash: failed to decode stdin: %w", err)
		}
		sources = append(sources, m)
	} else {
		for _, raw := range squashSources {
			var m map[string]string
			if err := json.Unmarshal([]byte(raw), &m); err != nil {
				return fmt.Errorf("squash: invalid source JSON: %w", err)
			}
			sources = append(sources, m)
		}
	}

	res, err := squash.Squash(sources, squash.Options{Strategy: strat})
	if err != nil {
		return err
	}

	if len(res.Conflicts) > 0 {
		fmt.Fprintf(cmd.ErrOrStderr(), "conflicts resolved: %v\n", res.Conflicts)
	}

	return json.NewEncoder(cmd.OutOrStdout()).Encode(res.Vars)
}
