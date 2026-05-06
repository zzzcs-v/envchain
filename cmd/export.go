package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/yourorg/envchain/internal/config"
	"github.com/yourorg/envchain/internal/context"
	"github.com/yourorg/envchain/internal/export"
)

var (
	exportFormat  string
	exportContext string
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Print resolved env vars for a context in the chosen format",
	Example: strings.Join([]string{
		"  envchain export --context staging",
		"  envchain export --context prod --format json",
		"  envchain export --context dev --format export >> .env",
	}, "\n"),
	RunE: runExport,
}

func init() {
	exportCmd.Flags().StringVarP(&exportContext, "context", "c", "", "context name to export (required)")
	exportCmd.Flags().StringVarP(&exportFormat, "format", "f", "dotenv", "output format: dotenv | export | json")
	_ = exportCmd.MarkFlagRequired("context")
	rootCmd.AddCommand(exportCmd)
}

func runExport(cmd *cobra.Command, _ []string) error {
	cfgPath, _ := cmd.Root().PersistentFlags().GetString("config")

	cfg, err := config.Load(cfgPath)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	resolver := context.NewResolver(cfg)
	env, err := resolver.Resolve(exportContext)
	if err != nil {
		return fmt.Errorf("resolving context %q: %w", exportContext, err)
	}

	fmt_, err := export.ParseFormat(exportFormat)
	if err != nil {
		return err
	}

	ex, err := export.NewExporter(fmt_)
	if err != nil {
		return err
	}

	return ex.Write(os.Stdout, env)
}
