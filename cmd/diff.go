package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/envchain/envchain/internal/config"
	"github.com/envchain/envchain/internal/context"
	"github.com/envchain/envchain/internal/diff"
)

var diffCmd = &cobra.Command{
	Use:   "diff <from-context> <to-context>",
	Short: "Show variable differences between two contexts",
	Args:  cobra.ExactArgs(2),
	RunE:  runDiff,
}

func init() {
	diffCmd.Flags().StringP("config", "c", "envchain.yaml", "path to config file")
	diffCmd.Flags().BoolP("show-unchanged", "u", false, "include unchanged variables in output")
	rootCmd.AddCommand(diffCmd)
}

func runDiff(cmd *cobra.Command, args []string) error {
	cfgPath, _ := cmd.Flags().GetString("config")
	showUnchanged, _ := cmd.Flags().GetBool("show-unchanged")

	fromCtx := args[0]
	toCtx := args[1]

	cfg, err := config.Load(cfgPath)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	resolver := context.NewResolver(cfg)

	fromVars, err := resolver.Resolve(fromCtx)
	if err != nil {
		return fmt.Errorf("resolving context %q: %w", fromCtx, err)
	}

	toVars, err := resolver.Resolve(toCtx)
	if err != nil {
		return fmt.Errorf("resolving context %q: %w", toCtx, err)
	}

	result := diff.Compare(fromVars, toVars)

	if !result.HasChanges() {
		fmt.Println("No differences found.")
		return nil
	}

	for _, c := range result.Changes {
		switch c.Type {
		case diff.Added:
			fmt.Fprintf(os.Stdout, "+ %s=%s\n", c.Key, c.NewValue)
		case diff.Removed:
			fmt.Fprintf(os.Stdout, "- %s=%s\n", c.Key, c.OldValue)
		case diff.Modified:
			fmt.Fprintf(os.Stdout, "~ %s: %s -> %s\n", c.Key, c.OldValue, c.NewValue)
		case diff.Unchanged:
			if showUnchanged {
				fmt.Fprintf(os.Stdout, "  %s=%s\n", c.Key, c.NewValue)
			}
		}
	}

	fmt.Println()
	fmt.Println(result.Summary())
	return nil
}
