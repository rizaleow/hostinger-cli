// Package describe emits the cobra command tree as machine-readable JSON, so
// that AI agents and scripts can discover the CLI surface without parsing
// `--help` text.
package describe

import (
	"encoding/json"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type cmdNode struct {
	Path        string     `json:"path"`
	Use         string     `json:"use"`
	Short       string     `json:"short,omitempty"`
	Long        string     `json:"long,omitempty"`
	Hidden      bool       `json:"hidden,omitempty"`
	Deprecated  string     `json:"deprecated,omitempty"`
	Flags       []flagNode `json:"flags,omitempty"`
	Subcommands []cmdNode  `json:"subcommands,omitempty"`
}

type flagNode struct {
	Name      string `json:"name"`
	Shorthand string `json:"shorthand,omitempty"`
	Type      string `json:"type"`
	Default   string `json:"default,omitempty"`
	Usage     string `json:"usage,omitempty"`
	Persist   bool   `json:"persistent,omitempty"`
}

// NewCmd builds the describe command, retaining a reference to the root so
// it can walk the whole tree even if invoked at a non-root path.
func NewCmd(root *cobra.Command) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "describe",
		Short: "Emit the command tree as JSON (machine-readable; for agents)",
		Long: "Prints the full command tree, including flags and usage, as a single JSON " +
			"document. Designed so that AI agents can discover the CLI surface programmatically.",
		Annotations: map[string]string{"skip-state": "true"},
		RunE: func(cmd *cobra.Command, args []string) error {
			tree := walk("", root)
			enc := json.NewEncoder(cmd.OutOrStdout())
			enc.SetIndent("", "  ")
			return enc.Encode(tree)
		},
	}
	return cmd
}

func walk(prefix string, c *cobra.Command) cmdNode {
	path := c.Name()
	if prefix != "" {
		path = prefix + " " + c.Name()
	}
	n := cmdNode{
		Path:       path,
		Use:        c.Use,
		Short:      c.Short,
		Long:       c.Long,
		Hidden:     c.Hidden,
		Deprecated: c.Deprecated,
	}
	collect := func(f *pflag.Flag, persist bool) {
		n.Flags = append(n.Flags, flagNode{
			Name:      f.Name,
			Shorthand: f.Shorthand,
			Type:      f.Value.Type(),
			Default:   f.DefValue,
			Usage:     f.Usage,
			Persist:   persist,
		})
	}
	c.PersistentFlags().VisitAll(func(f *pflag.Flag) { collect(f, true) })
	c.LocalNonPersistentFlags().VisitAll(func(f *pflag.Flag) { collect(f, false) })

	for _, sub := range c.Commands() {
		if sub.Hidden || sub.Name() == "help" {
			continue
		}
		n.Subcommands = append(n.Subcommands, walk(path, sub))
	}
	return n
}
