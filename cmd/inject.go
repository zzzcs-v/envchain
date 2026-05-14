package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/envchain/envchain/internal/inject"
	"github.com/spf13/cobra"
)

func init() {
	var (
		overwrite bool
		prefix    string
		sources   []string
	)

	cmd := &cobra.Command{
		Use:   "inject",
		Short: "Inject variables from named sources into a merged output",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runInject(sources, prefix, overwrite, cmd)
		},
	}

	cmd.Flags().BoolVar(&overwrite, "overwrite", false, "overwrite existing keys")
	cmd.Flags().StringVar(&prefix, "prefix", "", "prefix to prepend to all injected keys")
	cmd.Flags().StringArrayVar(&sources, "source", nil, "KEY=VALUE pairs as name:KEY=VALUE (repeatable)")

	rootCmd.AddCommand(cmd)
}

// runInject parses --source flags of the form "name:KEY=VALUE,..." and injects
// them into an empty destination map, printing JSON results.
func runInject(sources []string, prefix string, overwrite bool, cmd *cobra.Command) error {
	parsed, err := parseSources(sources)
	if err != nil {
		return err
	}

	dst := map[string]string{}
	inj := inject.New(inject.Options{Overwrite: overwrite, Prefix: prefix})
	results, err := inj.Inject(dst, parsed)
	if err != nil {
		return err
	}

	for _, r := range results {
		fmt.Fprintf(cmd.OutOrStdout(), "source=%s injected=%d skipped=%d\n", r.Source, r.Injected, r.Skipped)
	}

	return json.NewEncoder(os.Stdout).Encode(dst)
}

// parseSources converts "name:KEY=VALUE" strings into inject.Source slices.
func parseSources(raw []string) ([]inject.Source, error) {
	var out []inject.Source
	for _, s := range raw {
		colon := strings.IndexByte(s, ':')
		if colon < 1 {
			return nil, fmt.Errorf("inject: invalid source format %q, expected name:KEY=VALUE", s)
		}
		name := s[:colon]
		vars := map[string]string{}
		for _, pair := range strings.Split(s[colon+1:], ",") {
			eq := strings.IndexByte(pair, '=')
			if eq < 1 {
				return nil, fmt.Errorf("inject: invalid pair %q in source %q", pair, name)
			}
			vars[pair[:eq]] = pair[eq+1:]
		}
		out = append(out, inject.Source{Name: name, Vars: vars})
	}
	return out, nil
}
