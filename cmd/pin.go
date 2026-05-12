package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"text/tabwriter"

	"github.com/nicholasgasior/envchain/internal/pin"
	"github.com/spf13/cobra"
)

var pinDir string

func init() {
	pinCmd := &cobra.Command{
		Use:   "pin",
		Short: "Pin and retrieve named environment variable sets",
	}
	pinCmd.PersistentFlags().StringVar(&pinDir, "pin-dir", defaultPinDir(), "directory to store pins")

	pinCmd.AddCommand(
		&cobra.Command{
			Use:   "save <name> <context> KEY=VAL...",
			Short: "Save a pin",
			Args:  cobra.MinimumNArgs(2),
			RunE:  runPinSave,
		},
		&cobra.Command{
			Use:   "show <name>",
			Short: "Show a pin",
			Args:  cobra.ExactArgs(1),
			RunE:  runPinShow,
		},
		&cobra.Command{
			Use:   "list",
			Short: "List all pins",
			RunE:  runPinList,
		},
		&cobra.Command{
			Use:   "delete <name>",
			Short: "Delete a pin",
			Args:  cobra.ExactArgs(1),
			RunE:  runPinDelete,
		},
	)
	rootCmd.AddCommand(pinCmd)
}

func defaultPinDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".envchain", "pins")
}

func storeFromPinFlag() (*pin.Store, error) {
	return pin.NewStore(pinDir)
}

func runPinSave(cmd *cobra.Command, args []string) error {
	name, context := args[0], args[1]
	vars := make(map[string]string)
	for _, kv := range args[2:] {
		for i, c := range kv {
			if c == '=' {
				vars[kv[:i]] = kv[i+1:]
				break
			}
		}
	}
	s, err := storeFromPinFlag()
	if err != nil {
		return err
	}
	if err := s.Save(pin.Pin{Name: name, Context: context, Vars: vars}); err != nil {
		return err
	}
	fmt.Fprintf(cmd.OutOrStdout(), "pinned %q (%s)\n", name, context)
	return nil
}

func runPinShow(cmd *cobra.Command, args []string) error {
	s, err := storeFromPinFlag()
	if err != nil {
		return err
	}
	p, err := s.Load(args[0])
	if err != nil {
		return err
	}
	w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
	fmt.Fprintf(w, "name:\t%s\ncontext:\t%s\ncreated:\t%s\n", p.Name, p.Context, p.CreatedAt.Format("2006-01-02 15:04:05"))
	for k, v := range p.Vars {
		fmt.Fprintf(w, "%s\t%s\n", k, v)
	}
	return w.Flush()
}

func runPinList(cmd *cobra.Command, _ []string) error {
	s, err := storeFromPinFlag()
	if err != nil {
		return err
	}
	pins, err := s.List()
	if err != nil {
		return err
	}
	if len(pins) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "no pins saved")
		return nil
	}
	w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "NAME\tCONTEXT\tCREATED")
	for _, p := range pins {
		fmt.Fprintf(w, "%s\t%s\t%s\n", p.Name, p.Context, p.CreatedAt.Format("2006-01-02 15:04:05"))
	}
	return w.Flush()
}

func runPinDelete(cmd *cobra.Command, args []string) error {
	s, err := storeFromPinFlag()
	if err != nil {
		return err
	}
	if err := s.Delete(args[0]); err != nil {
		return err
	}
	fmt.Fprintf(cmd.OutOrStdout(), "deleted pin %q\n", args[0])
	return nil
}
