package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/envchain/envchain/internal/history"
)

var historyDir string

func init() {
	historyCmd.Flags().StringVar(&historyDir, "dir", defaultHistoryDir(), "directory to store history entries")
	historyCmd.AddCommand(historyClearCmd)
	rootCmd.AddCommand(historyCmd)
}

func defaultHistoryDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ".envchain/history"
	}
	return home + "/.envchain/history"
}

var historyCmd = &cobra.Command{
	Use:   "history",
	Short: "Show past context resolutions",
	RunE:  runHistory,
}

var historyClearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear all history entries",
	RunE:  runHistoryClear,
}

func runHistory(cmd *cobra.Command, args []string) error {
	s, err := history.NewStore(historyDir)
	if err != nil {
		return fmt.Errorf("history: %w", err)
	}
	entries, err := s.List()
	if err != nil {
		return fmt.Errorf("history: %w", err)
	}
	if len(entries) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "no history entries found")
		return nil
	}
	w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "TIMESTAMP\tCONTEXT\tFORMAT\tVARS")
	for _, e := range entries {
		fmt.Fprintf(w, "%s\t%s\t%s\t%d\n",
			e.Timestamp.Format("2006-01-02 15:04:05"),
			e.Context,
			e.Format,
			len(e.Vars),
		)
	}
	return w.Flush()
}

func runHistoryClear(cmd *cobra.Command, args []string) error {
	s, err := history.NewStore(historyDir)
	if err != nil {
		return fmt.Errorf("history clear: %w", err)
	}
	if err := s.Clear(); err != nil {
		return fmt.Errorf("history clear: %w", err)
	}
	fmt.Fprintln(cmd.OutOrStdout(), "history cleared")
	return nil
}
