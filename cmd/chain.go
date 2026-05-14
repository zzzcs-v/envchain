package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"envchain/internal/chain"
)

var chainStore = chain.New()

func init() {
	chainCmd := &cobra.Command{
		Use:   "chain",
		Short: "Manage named context chains",
	}

	setCmd := &cobra.Command{
		Use:   "set <name> <ctx1,ctx2,...>",
		Short: "Define a named chain of contexts",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctxs := strings.Split(args[1], ",")
			return chainStore.Set(args[0], ctxs)
		},
	}

	getCmd := &cobra.Command{
		Use:   "show <name>",
		Short: "Show contexts in a named chain",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := chainStore.Get(args[0])
			if err != nil {
				return err
			}
			fmt.Fprintf(os.Stdout, "chain: %s\n", c.Name)
			for i, ctx := range c.Contexts {
				fmt.Fprintf(os.Stdout, "  %d. %s\n", i+1, ctx)
			}
			return nil
		},
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all named chains",
		RunE: func(cmd *cobra.Command, args []string) error {
			chains := chainStore.List()
			if len(chains) == 0 {
				fmt.Fprintln(os.Stdout, "no chains defined")
				return nil
			}
			for _, c := range chains {
				fmt.Fprintf(os.Stdout, "%s: %s\n", c.Name, strings.Join(c.Contexts, " -> "))
			}
			return nil
		},
	}

	delCmd := &cobra.Command{
		Use:   "delete <name>",
		Short: "Delete a named chain",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return chainStore.Delete(args[0])
		},
	}

	chainCmd.AddCommand(setCmd, getCmd, listCmd, delCmd)
	rootCmd.AddCommand(chainCmd)
}
