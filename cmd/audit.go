package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"

	"github.com/yourorg/envchain/internal/audit"
)

var (
	auditLogFile string
	auditLimit   int
)

func init() {
	auditCmd := &cobra.Command{
		Use:   "audit",
		Short: "Show audit log of envchain actions",
		RunE:  runAudit,
	}
	auditCmd.Flags().StringVar(&auditLogFile, "log", ".envchain-audit.log", "path to audit log file")
	auditCmd.Flags().IntVar(&auditLimit, "limit", 20, "max number of entries to show (0 = all)")
	rootCmd.AddCommand(auditCmd)
}

func runAudit(cmd *cobra.Command, args []string) error {
	entries, err := audit.ReadAll(auditLogFile)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Fprintln(cmd.OutOrStdout(), "no audit log found")
			return nil
		}
		return fmt.Errorf("reading audit log: %w", err)
	}

	if len(entries) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "audit log is empty")
		return nil
	}

	// show most recent first
	start := 0
	if auditLimit > 0 && len(entries) > auditLimit {
		start = len(entries) - auditLimit
	}
	visible := entries[start:]

	w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "TIMESTAMP\tACTION\tCONTEXT\tVARS")
	for i := len(visible) - 1; i >= 0; i-- {
		e := visible[i]
		ts := e.Timestamp.Format(time.RFC3339)
		varsStr := fmt.Sprintf("%d", len(e.Vars))
		if len(e.Vars) == 0 {
			varsStr = "-"
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", ts, e.Action, e.Context, varsStr)
	}
	return w.Flush()
}
