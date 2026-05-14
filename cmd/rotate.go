package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/user/envchain/internal/rotate"
)

var rotateCmd = &cobra.Command{
	Use:   "rotate",
	Short: "Rotate values for keys in a context",
	RunE:  runRotate,
}

func init() {
	rotateCmd.Flags().StringP("context", "c", "default", "context name to rotate keys in")
	rotateCmd.Flags().StringSliceP("keys", "k", nil, "keys to rotate (default: all)")
	rotateCmd.Flags().BoolP("overwrite", "o", false, "overwrite non-empty values")
	rotateCmd.Flags().StringP("suffix", "s", "rotated", "suffix appended to old value to generate new value")
	RootCmd.AddCommand(rotateCmd)
}

func runRotate(cmd *cobra.Command, _ []string) error {
	ctxName, _ := cmd.Flags().GetString("context")
	keys, _ := cmd.Flags().GetStringSlice("keys")
	overwrite, _ := cmd.Flags().GetBool("overwrite")
	suffix, _ := cmd.Flags().GetString("suffix")

	// Minimal stub vars — in real usage these come from the loaded config/context.
	vars := map[string]string{
		"API_KEY":    os.Getenv("API_KEY"),
		"DB_PASS":    os.Getenv("DB_PASS"),
		"SECRET_TOKEN": os.Getenv("SECRET_TOKEN"),
	}

	gen := func(key string) (string, error) {
		old := vars[key]
		if old == "" {
			return strings.ToLower(key) + "_" + suffix, nil
		}
		return old + "_" + suffix, nil
	}

	res, err := rotate.Rotate(vars, rotate.Options{
		Context:   ctxName,
		Keys:      keys,
		Overwrite: overwrite,
		Generator: gen,
	})
	if err != nil {
		return fmt.Errorf("rotate: %w", err)
	}

	enc := json.NewEncoder(cmd.OutOrStdout())
	enc.SetIndent("", "  ")
	if err := enc.Encode(res.Rotated); err != nil {
		return err
	}

	fmt.Fprintln(cmd.OutOrStdout(), res.Summary())

	if len(res.Errors) > 0 {
		for _, e := range res.Errors {
			fmt.Fprintln(cmd.ErrOrStderr(), "error:", e)
		}
		return fmt.Errorf("rotation completed with %d error(s)", len(res.Errors))
	}
	return nil
}
