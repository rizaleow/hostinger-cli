package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/rizaleow/hostinger-cli/internal/clictx"
	"github.com/rizaleow/hostinger-cli/internal/config"
	"github.com/rizaleow/hostinger-cli/internal/output"
	"github.com/rizaleow/hostinger-cli/internal/version"
)

type rootFlags struct {
	configPath string
	profile    string
	token      string
	baseURL    string
	outputFmt  string
	template   string
	jq         string
	noColor    bool
	compact    bool
	verbose    bool
}

// Execute runs the root command with ctx (cancelled on SIGINT/SIGTERM by main).
func Execute(ctx context.Context) error {
	return newRootCmd().ExecuteContext(ctx)
}

func newRootCmd() *cobra.Command {
	rf := &rootFlags{}
	root := &cobra.Command{
		Use:   "hostinger-cli",
		Short: "Command-line interface for the Hostinger API",
		Long: "hostinger-cli is a Go-native CLI for managing Hostinger resources: VPS, " +
			"DNS, domains, hosting, billing, and email marketing.\n\n" +
			"Authentication: set HOSTINGER_API_TOKEN, run `hostinger-cli auth login`, " +
			"or pass --token. JSON is emitted automatically when stdout is piped; pass " +
			"--output table to force a human-readable table.",
		Version:           version.String(),
		SilenceUsage:      true,
		SilenceErrors:     true,
		DisableAutoGenTag: true,
	}

	root.PersistentFlags().StringVar(&rf.configPath, "config", "", "config file (default ~/.config/hostinger-cli/config.yaml)")
	root.PersistentFlags().StringVar(&rf.profile, "profile", "", "config profile to use")
	root.PersistentFlags().StringVar(&rf.token, "token", "", "API token (overrides config and env)")
	root.PersistentFlags().StringVar(&rf.baseURL, "base-url", "", "override Hostinger API base URL")
	root.PersistentFlags().StringVarP(&rf.outputFmt, "output", "o", "", "output format: table|json|yaml|template (auto-detected from TTY by default)")
	root.PersistentFlags().StringVar(&rf.template, "template", "", "Go text/template body (with --output template)")
	root.PersistentFlags().StringVar(&rf.jq, "jq", "", "filter JSON output through this gojq expression")
	root.PersistentFlags().BoolVar(&rf.noColor, "no-color", false, "disable colored output")
	root.PersistentFlags().BoolVar(&rf.compact, "compact", false, "compact JSON output (no indentation)")
	root.PersistentFlags().BoolVarP(&rf.verbose, "verbose", "v", false, "verbose logging")

	root.SetVersionTemplate("{{.Version}}\n")

	root.PersistentPreRunE = func(cmd *cobra.Command, _ []string) error {
		// Commands that don't need state can skip prerun by adding the
		// annotation "skip-state":"true".
		if cmd.Annotations["skip-state"] == "true" {
			return nil
		}

		path := rf.configPath
		if path == "" {
			p, err := config.DefaultPath()
			if err != nil {
				return err
			}
			path = p
		}
		file, err := config.Load(path)
		if err != nil {
			return err
		}
		resolved := config.Resolve(file, config.ResolveOptions{
			FlagToken:   rf.token,
			FlagBaseURL: rf.baseURL,
			Profile:     rf.profile,
		})

		state := &clictx.State{
			Config:     file,
			ConfigPath: path,
			Resolved:   resolved,
			OutputOptions: output.Options{
				Format:   output.Format(rf.outputFmt),
				Template: rf.template,
				JQ:       rf.jq,
				NoColor:  rf.noColor || os.Getenv("NO_COLOR") != "",
				Compact:  rf.compact,
				Out:      cmd.OutOrStdout(),
			},
		}
		cmd.SetContext(clictx.With(cmd.Context(), state))
		return nil
	}

	registerCommands(root)
	return root
}

// ErrUsage wraps an error with the cobra usage banner.
func ErrUsage(cmd *cobra.Command, err error) error {
	cmd.SilenceUsage = false
	return fmt.Errorf("%w", err)
}
