package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/envchain/envchain/internal/group"
	"github.com/spf13/cobra"
)

func defaultGroupDir() string {
	home, _ := os.UserHomeDir()
	return home + "/.envchain/groups"
}

func storeFromGroupFlag(dir string) (*group.Store, error) {
	if dir == "" {
		dir = defaultGroupDir()
	}
	return group.NewStore(dir)
}

func init() {
	var dir string

	groupCmd := &cobra.Command{
		Use:   "group",
		Short: "Manage named groups of contexts",
	}
	groupCmd.PersistentFlags().StringVar(&dir, "dir", "", "group storage directory")

	setCmd := &cobra.Command{
		Use:   "set <name> <ctx1,ctx2,...>",
		Short: "Create or update a group",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := storeFromGroupFlag(dir)
			if err != nil {
				return err
			}
			ctxs := strings.Split(args[1], ",")
			return s.Save(group.Group{Name: args[0], Contexts: ctxs})
		},
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all groups",
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := storeFromGroupFlag(dir)
			if err != nil {
				return err
			}
			groups, err := s.List()
			if err != nil {
				return err
			}
			if len(groups) == 0 {
				fmt.Println("no groups defined")
				return nil
			}
			for _, g := range groups {
				fmt.Printf("%s\t[%s]\n", g.Name, strings.Join(g.Contexts, ", "))
			}
			return nil
		},
	}

	deleteCmd := &cobra.Command{
		Use:   "delete <name>",
		Short: "Delete a group",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := storeFromGroupFlag(dir)
			if err != nil {
				return err
			}
			return s.Delete(args[0])
		},
	}

	groupCmd.AddCommand(setCmd, listCmd, deleteCmd)
	rootCmd.AddCommand(groupCmd)
}
