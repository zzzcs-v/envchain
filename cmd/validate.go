package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yourorg/envchain/internal/config"
	"github.com/yourorg/envchain/internal/context"
	"github.com/yourorg/envchain/internal/validate"
)

var validateCmd = &cobra.Command{
	Use:   "validate [context]",
	Short: "Validate env var keys and values for a given context",
	Args:  cobra.ExactArgs(1),
	RunE:  runValidate,
}

func init() {
	rootCmd.AddCommand(validateCmd)
	validateCmd.Flags().StringP("config", "c", "envchain.yaml", "path to config file")
}

func runValidate(cmd *cobra.Command, args []string) error {
	cfgPath, _ := cmd.Flags().GetString("config")
	ctxName := args[0]

	cfg, err := config.Load(cfgPath)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	resolver := context.NewResolver(cfg)
	vars, err := resolver.Resolve(ctxName)
	if err != nil {
		return fmt.Errorf("resolving context %q: %w", ctxName, err)
	}

	result := validate.Vars(vars)
	if summary := result.Summary(); summary != "" {
		fmt.Fprint(os.Stderr, summary)
	}

	if !result.OK() {
		return fmt.Errorf("validation failed with %d error(s)", len(result.Errors))
	}

	fmt.Printf("context %q passed validation (%d var(s) checked, %d warning(s))\n",
		ctxName, len(vars), len(result.Warnings))
	return nil
}
