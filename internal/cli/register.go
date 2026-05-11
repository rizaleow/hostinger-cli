package cli

import (
	"github.com/spf13/cobra"

	"github.com/rizaleow/hostinger-cli/internal/cli/auth"
	"github.com/rizaleow/hostinger-cli/internal/cli/billing"
	configcmd "github.com/rizaleow/hostinger-cli/internal/cli/config"
	"github.com/rizaleow/hostinger-cli/internal/cli/describe"
	"github.com/rizaleow/hostinger-cli/internal/cli/dns"
	"github.com/rizaleow/hostinger-cli/internal/cli/domains"
	"github.com/rizaleow/hostinger-cli/internal/cli/hosting"
	"github.com/rizaleow/hostinger-cli/internal/cli/reach"
	"github.com/rizaleow/hostinger-cli/internal/cli/vps"
)

func registerCommands(root *cobra.Command) {
	root.AddCommand(
		auth.NewCmd(),
		configcmd.NewCmd(),
		describe.NewCmd(root),
		newCompletionCmd(),
		billing.NewCmd(),
		dns.NewCmd(),
		domains.NewCmd(),
		hosting.NewCmd(),
		reach.NewCmd(),
		vps.NewCmd(),
	)
}
