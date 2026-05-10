package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envchain/internal/config"
	"envchain/internal/prompt"
)

var promptCmd = &cobra.Command{
	Use:   "prompt",
	Short: "Interactively set env vars for a context",
	RunE:  runPrompt,
}

var promptContext string
var promptConfig string

func init() {
	promptCmd.Flags().StringVarP(&promptContext, "context", "c", "", "context to populate (required)")
	promptCmd.Flags().StringVarP(&promptConfig, "config", "f", "envchain.yaml", "path to config file")
	_ = promptCmd.MarkFlagRequired("context")
	rootCmd.AddCommand(promptCmd)
}

func runPrompt(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load(promptConfig)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	var ctx *config.Context
	for i := range cfg.Contexts {
		if cfg.Contexts[i].Name == promptContext {
			ctx = &cfg.Contexts[i]
			break
		}
	}
	if ctx == nil {
		return fmt.Errorf("context %q not found", promptContext)
	}

	p := prompt.New()

	ok, err := p.Confirm(fmt.Sprintf("Populate %d var(s) for context %q", len(ctx.Vars), promptContext))
	if err != nil {
		return err
	}
	if !ok {
		fmt.Fprintln(os.Stderr, "aborted")
		return nil
	}

	for key := range ctx.Vars {
		val, err := p.Ask(fmt.Sprintf("  %s", key))
		if err != nil {
			return fmt.Errorf("prompt for %q: %w", key, err)
		}
		ctx.Vars[key] = val
	}

	fmt.Fprintf(os.Stdout, "\n✓ %d variable(s) collected for context %q\n", len(ctx.Vars), promptContext)
	return nil
}
