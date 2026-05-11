// Package domains implements `hostinger-cli domains`.
package domains

import "github.com/spf13/cobra"

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "domains", Short: "Domain portfolio, WHOIS, forwarding, availability"}
	cmd.AddCommand(
		newPortfolioCmd(),
		newAvailabilityCmd(),
		newForwardingCmd(),
		newWHOISCmd(),
		newVerifierCmd(),
	)
	return cmd
}
