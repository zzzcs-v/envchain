package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/envchain/envchain/internal/sort"
)

var (
	sortBy    string
	sortOrder string
)

func init() {
	sortCmd := &cobra.Command{
		Use:   "sort",
		Short: "Sort environment variables by key or value",
		Long:  "Reads a JSON env map from stdin and outputs a sorted JSON array of {key, value} pairs.",
		RunE:  runSort,
	}
	sortCmd.Flags().StringVar(&sortBy, "by", "key", "Sort field: key or value")
	sortCmd.Flags().StringVar(&sortOrder, "order", "asc", "Sort order: asc or desc")
	rootCmd.AddCommand(sortCmd)
}

func runSort(cmd *cobra.Command, args []string) error {
	var src map[string]string
	if err := json.NewDecoder(os.Stdin).Decode(&src); err != nil {
		return fmt.Errorf("failed to decode stdin: %w", err)
	}

	res, err := sort.Sort(src, sort.Options{
		By:    sortBy,
		Order: sort.Order(sortOrder),
	})
	if err != nil {
		return err
	}

	out, err := json.MarshalIndent(res.Pairs, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to encode output: %w", err)
	}

	fmt.Fprintln(cmd.OutOrStdout(), string(out))
	return nil
}
