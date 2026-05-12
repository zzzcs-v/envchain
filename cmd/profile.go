package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"envchain/internal/profile"
)

var (
	profileDir     string
	profileContext string
	profileFormat  string
)

func init() {
	profileCmd := &cobra.Command{
		Use:   "profile",
		Short: "Manage named profiles (context + format presets)",
	}

	saveCmd := &cobra.Command{
		Use:   "save <name>",
		Short: "Save a profile",
		Args:  cobra.ExactArgs(1),
		RunE:  runProfileSave,
	}
	saveCmd.Flags().StringVar(&profileContext, "context", "", "context name (required)")
	saveCmd.Flags().StringVar(&profileFormat, "format", "dotenv", "export format")
	_ = saveCmd.MarkFlagRequired("context")

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List saved profiles",
		RunE:  runProfileList,
	}

	showCmd := &cobra.Command{
		Use:   "show <name>",
		Short: "Show a profile",
		Args:  cobra.ExactArgs(1),
		RunE:  runProfileShow,
	}

	deleteCmd := &cobra.Command{
		Use:   "delete <name>",
		Short: "Delete a profile",
		Args:  cobra.ExactArgs(1),
		RunE:  runProfileDelete,
	}

	profileCmd.PersistentFlags().StringVar(&profileDir, "profile-dir", ".envchain/profiles", "directory for profiles")
	profileCmd.AddCommand(saveCmd, listCmd, showCmd, deleteCmd)
	rootCmd.AddCommand(profileCmd)
}

func storeFromFlag() (*profile.Store, error) {
	return profile.NewStore(profileDir)
}

func runProfileSave(cmd *cobra.Command, args []string) error {
	s, err := storeFromFlag()
	if err != nil {
		return err
	}
	p := profile.Profile{Name: args[0], Context: profileContext, Format: profileFormat}
	if err := s.Save(p); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "profile %q saved\n", args[0])
	return nil
}

func runProfileList(cmd *cobra.Command, args []string) error {
	s, err := storeFromFlag()
	if err != nil {
		return err
	}
	names, err := s.List()
	if err != nil {
		return err
	}
	if len(names) == 0 {
		fmt.Fprintln(os.Stdout, "no profiles saved")
		return nil
	}
	fmt.Fprintln(os.Stdout, strings.Join(names, "\n"))
	return nil
}

func runProfileShow(cmd *cobra.Command, args []string) error {
	s, err := storeFromFlag()
	if err != nil {
		return err
	}
	p, err := s.Load(args[0])
	if err != nil {
		return err
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(p)
}

func runProfileDelete(cmd *cobra.Command, args []string) error {
	s, err := storeFromFlag()
	if err != nil {
		return err
	}
	if err := s.Delete(args[0]); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "profile %q deleted\n", args[0])
	return nil
}
