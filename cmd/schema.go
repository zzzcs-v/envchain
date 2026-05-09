package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/user/envchain/internal/config"
	"github.com/user/envchain/internal/context"
	"github.com/user/envchain/internal/schema"
)

var (
	schemaFile    string
	schemaContext string
)

func init() {
	schemaCmd := &cobra.Command{
		Use:   "schema",
		Short: "Validate resolved context vars against a JSON schema file",
		RunE:  runSchema,
	}
	schemaCmd.Flags().StringVarP(&schemaFile, "schema", "s", "", "path to JSON schema file (required)")
	schemaCmd.Flags().StringVarP(&schemaContext, "context", "c", "", "context name to validate (required)")
	_ = schemaCmd.MarkFlagRequired("schema")
	_ = schemaCmd.MarkFlagRequired("context")
	rootCmd.AddCommand(schemaCmd)
}

func runSchema(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load(cfgFile)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	resolver := context.NewResolver(cfg)
	vars, err := resolver.Resolve(schemaContext)
	if err != nil {
		return fmt.Errorf("resolve context %q: %w", schemaContext, err)
	}

	data, err := os.ReadFile(schemaFile)
	if err != nil {
		return fmt.Errorf("read schema file: %w", err)
	}

	var rawRules map[string]schema.Rule
	if err := json.Unmarshal(data, &rawRules); err != nil {
		return fmt.Errorf("parse schema file: %w", err)
	}

	s, err := schema.NewSchema(rawRules)
	if err != nil {
		return fmt.Errorf("build schema: %w", err)
	}

	violations := s.Validate(vars)
	if len(violations) == 0 {
		fmt.Fprintf(cmd.OutOrStdout(), "✓ context %q passes schema validation\n", schemaContext)
		return nil
	}

	fmt.Fprintf(cmd.OutOrStdout(), "✗ schema violations in context %q:\n", schemaContext)
	for _, v := range violations {
		fmt.Fprintf(cmd.OutOrStdout(), "  - %s\n", v)
	}
	return fmt.Errorf("%d violation(s) found", len(violations))
}
