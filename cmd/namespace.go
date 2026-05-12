package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"envchain/internal/namespace"
)

func defaultNamespaceDir() string {
	home, _ := os.UserHomeDir()
	return home + "/.envchain/namespaces"
}

func storeFromNSFlag(dir string) (*namespace.Store, error) {
	return namespace.NewStore(dir)
}

func init() {
	var dir string

	nsCmd := &cobra.Command{
		Use:   "namespace",
		Short: "Manage environment namespaces",
	}

	saveCmd := &cobra.Command{
		Use:   "save <name> <prefix> <ctx,...>",
		Short: "Save a namespace entry",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := storeFromNSFlag(dir)
			if err != nil {
				return err
			}
			ctxs := strings.Split(args[2], ",")
			return s.Save(namespace.Entry{Name: args[0], Prefix: args[1], Contexts: ctxs})
		},
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all namespaces",
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := storeFromNSFlag(dir)
			if err != nil {
				return err
			}
			entries, err := s.List()
			if err != nil {
				return err
			}
			if len(entries) == 0 {
				fmt.Println("no namespaces defined")
				return nil
			}
			for _, e := range entries {
				fmt.Printf("%-20s prefix=%-15s contexts=%s\n", e.Name, e.Prefix, strings.Join(e.Contexts, ","))
			}
			return nil
		},
	}

	delCmd := &cobra.Command{
		Use:   "delete <name>",
		Short: "Delete a namespace",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := storeFromNSFlag(dir)
			if err != nil {
				return err
			}
			return s.Delete(args[0])
		},
	}

	for _, sub := range []*cobra.Command{saveCmd, listCmd, delCmd} {
		sub.Flags().StringVar(&dir, "dir", defaultNamespaceDir(), "namespace storage directory")
		nsCmd.AddCommand(sub)
	}
	rootCmd.AddCommand(nsCmd)
}
