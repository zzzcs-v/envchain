package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/user/envchain/internal/tag"
)

var tagFile string

func defaultTagDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".envchain", "tags.json")
}

func init() {
	tagCmd := &cobra.Command{
		Use:   "tag",
		Short: "Manage context tags",
	}
	tagCmd.PersistentFlags().StringVar(&tagFile, "store", defaultTagDir(), "path to tag store")

	setCmd := &cobra.Command{
		Use:   "set <name> <context,...>",
		Short: "Create or update a tag with associated contexts",
		Args:  cobra.ExactArgs(2),
		RunE:  runTagSet,
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all tags",
		RunE:  runTagList,
	}

	deleteCmd := &cobra.Command{
		Use:   "delete <name>",
		Short: "Delete a tag",
		Args:  cobra.ExactArgs(1),
		RunE:  runTagDelete,
	}

	tagCmd.AddCommand(setCmd, listCmd, deleteCmd)
	rootCmd.AddCommand(tagCmd)
}

func storeFromTagFlag() (*tag.Store, error) {
	return tag.NewStore(tagFile)
}

func runTagSet(cmd *cobra.Command, args []string) error {
	name := args[0]
	contexts := strings.Split(args[1], ",")
	s, err := storeFromTagFlag()
	if err != nil {
		return err
	}
	if err := s.Set(name, contexts); err != nil {
		return err
	}
	fmt.Fprintf(cmd.OutOrStdout(), "tag %q saved with %d context(s)\n", name, len(contexts))
	return nil
}

func runTagList(cmd *cobra.Command, _ []string) error {
	s, err := storeFromTagFlag()
	if err != nil {
		return err
	}
	tags, err := s.List()
	if err != nil {
		return err
	}
	if len(tags) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "no tags found")
		return nil
	}
	for _, t := range tags {
		fmt.Fprintf(cmd.OutOrStdout(), "%-20s %s\n", t.Name, strings.Join(t.Contexts, ", "))
	}
	return nil
}

func runTagDelete(cmd *cobra.Command, args []string) error {
	s, err := storeFromTagFlag()
	if err != nil {
		return err
	}
	if err := s.Delete(args[0]); err != nil {
		return err
	}
	fmt.Fprintf(cmd.OutOrStdout(), "tag %q deleted\n", args[0])
	return nil
}
