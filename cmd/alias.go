package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/envchain/envchain/internal/alias"
)

func defaultAliasDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".envchain", "aliases")
}

func storeFromAliasFlag(dir string) (*alias.Store, error) {
	if dir == "" {
		dir = defaultAliasDir()
	}
	return alias.NewStore(dir)
}

func init() {
	aliasCmd := &cobra.Command{
		Use:   "alias",
		Short: "Manage context aliases",
	}

	var aliasDir string

	setCmd := &cobra.Command{
		Use:   "set <name> <context>",
		Short: "Create or update an alias",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAliasSet(aliasDir, args[0], args[1])
		},
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all aliases",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAliasList(aliasDir)
		},
	}

	deleteCmd := &cobra.Command{
		Use:   "delete <name>",
		Short: "Delete an alias",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAliasDelete(aliasDir, args[0])
		},
	}

	for _, sub := range []*cobra.Command{setCmd, listCmd, deleteCmd} {
		sub.Flags().StringVar(&aliasDir, "alias-dir", "", "directory for alias storage")
		aliasCmd.AddCommand(sub)
	}
	rootCmd.AddCommand(aliasCmd)
}

func runAliasSet(dir, name, context string) error {
	s, err := storeFromAliasFlag(dir)
	if err != nil {
		return err
	}
	if err := s.Set(name, context); err != nil {
		return err
	}
	fmt.Printf("alias %q -> %q saved\n", name, context)
	return nil
}

func runAliasList(dir string) error {
	s, err := storeFromAliasFlag(dir)
	if err != nil {
		return err
	}
	aliases, err := s.List()
	if err != nil {
		return err
	}
	if len(aliases) == 0 {
		fmt.Println("no aliases defined")
		return nil
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "NAME\tCONTEXT")
	for _, a := range aliases {
		fmt.Fprintf(w, "%s\t%s\n", a.Name, a.Context)
	}
	return w.Flush()
}

func runAliasDelete(dir, name string) error {
	s, err := storeFromAliasFlag(dir)
	if err != nil {
		return err
	}
	if err := s.Delete(name); err != nil {
		return err
	}
	fmt.Printf("alias %q deleted\n", name)
	return nil
}
