package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/user/envchain/internal/scope"
)

func defaultScopeDir() string {
	home, _ := os.UserHomeDir()
	return home + "/.envchain/scopes"
}

func storeFromScopeFlag(dir string) (*scope.Store, error) {
	if dir == "" {
		dir = defaultScopeDir()
	}
	return scope.NewStore(dir)
}

func init() {
	var dir string

	scopeCmd := &cobra.Command{
		Use:   "scope",
		Short: "Manage environment variable scopes",
	}

	saveCmd := &cobra.Command{
		Use:   "save <name> KEY=VALUE...",
		Short: "Save a scope with the given key=value pairs",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			st, err := storeFromScopeFlag(dir)
			if err != nil {
				return err
			}
			vars := make(map[string]string)
			for _, pair := range args[1:] {
				parts := strings.SplitN(pair, "=", 2)
				if len(parts) != 2 {
					return fmt.Errorf("invalid pair %q: expected KEY=VALUE", pair)
				}
				vars[parts[0]] = parts[1]
			}
			return st.Save(scope.Scope{Name: args[0], Vars: vars})
		},
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all saved scopes",
		RunE: func(cmd *cobra.Command, args []string) error {
			st, err := storeFromScopeFlag(dir)
			if err != nil {
				return err
			}
			names, err := st.List()
			if err != nil {
				return err
			}
			for _, n := range names {
				fmt.Fprintln(cmd.OutOrStdout(), n)
			}
			return nil
		},
	}

	showCmd := &cobra.Command{
		Use:   "show <name>",
		Short: "Show variables in a scope as JSON",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			st, err := storeFromScopeFlag(dir)
			if err != nil {
				return err
			}
			sc, err := st.Load(args[0])
			if err != nil {
				return err
			}
			enc := json.NewEncoder(cmd.OutOrStdout())
			enc.SetIndent("", "  ")
			return enc.Encode(sc)
		},
	}

	scopeCmd.PersistentFlags().StringVar(&dir, "dir", "", "directory for scope storage")
	scopeCmd.AddCommand(saveCmd, listCmd, showCmd)
	rootCmd.AddCommand(scopeCmd)
}
