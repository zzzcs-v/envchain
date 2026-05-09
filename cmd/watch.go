package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"

	"github.com/envchain/envchain/internal/watch"
)

var (
	watchInterval string
)

func init() {
	watchCmd := &cobra.Command{
		Use:   "watch [config-file...]",
		Short: "Watch config files and print a notice when they change",
		Args:  cobra.MinimumNArgs(1),
		RunE:  runWatch,
	}
	watchCmd.Flags().StringVarP(&watchInterval, "interval", "i", "2s", "polling interval (e.g. 500ms, 2s)")
	rootCmd.AddCommand(watchCmd)
}

func runWatch(cmd *cobra.Command, args []string) error {
	dur, err := time.ParseDuration(watchInterval)
	if err != nil {
		return fmt.Errorf("invalid interval %q: %w", watchInterval, err)
	}

	w, err := watch.New(args, dur)
	if err != nil {
		return err
	}

	w.Start()
	defer w.Stop()

	fmt.Fprintf(cmd.OutOrStdout(), "watching %d file(s), interval %s — press Ctrl+C to stop\n", len(args), dur)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case ev := <-w.Events:
			fmt.Fprintf(cmd.OutOrStdout(), "changed  %s\n  old: %s\n  new: %s\n  at:  %s\n",
				ev.Path, ev.OldHash, ev.NewHash, ev.At.Format(time.RFC3339))
		case watchErr := <-w.Errors:
			fmt.Fprintf(cmd.ErrOrStderr(), "watch error: %v\n", watchErr)
		case <-sig:
			fmt.Fprintln(cmd.OutOrStdout(), "stopping watcher")
			return nil
		}
	}
}
