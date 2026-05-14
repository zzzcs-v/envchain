package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envchain/internal/compare"
)

var compareJSON bool

func init() {
	cmd := &cobra.Command{
		Use:   "compare <left-context> <right-context>",
		Short: "Compare two environment contexts and show differences",
		Args:  cobra.ExactArgs(2),
		RunE:  runCompare,
	}
	cmd.Flags().BoolVar(&compareJSON, "json", false, "output result as JSON")
	rootCmd.AddCommand(cmd)
}

func runCompare(cmd *cobra.Command, args []string) error {
	leftName, rightName := args[0], args[1]

	cfg, err := loadConfig(cmd)
	if err != nil {
		return err
	}

	left, ok := cfg.Contexts[leftName]
	if !ok {
		return fmt.Errorf("context %q not found", leftName)
	}
	right, ok := cfg.Contexts[rightName]
	if !ok {
		return fmt.Errorf("context %q not found", rightName)
	}

	res := compare.Compare(left.Vars, right.Vars)

	if compareJSON {
		return json.NewEncoder(os.Stdout).Encode(map[string]interface{}{
			"only_in_left":  res.OnlyInLeft,
			"only_in_right": res.OnlyInRight,
			"different":     res.Different,
			"shared":        res.Shared,
		})
	}

	fmt.Fprintf(os.Stdout, "Comparing %q vs %q\n", leftName, rightName)
	fmt.Fprintln(os.Stdout, res.Summary())

	for _, k := range res.Keys() {
		if v, ok := res.OnlyInLeft[k]; ok {
			fmt.Fprintf(os.Stdout, "  - %s=%s (only in %s)\n", k, v, leftName)
		} else if v, ok := res.OnlyInRight[k]; ok {
			fmt.Fprintf(os.Stdout, "  + %s=%s (only in %s)\n", k, v, rightName)
		} else if pair, ok := res.Different[k]; ok {
			fmt.Fprintf(os.Stdout, "  ~ %s: %s -> %s\n", k, pair[0], pair[1])
		}
	}

	return nil
}
