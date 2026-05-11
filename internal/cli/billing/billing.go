// Package billing implements the `hostinger-cli billing` command tree.
package billing

import "github.com/spf13/cobra"

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "billing", Short: "Billing: catalog, payment methods, subscriptions"}
	cmd.AddCommand(newCatalogCmd(), newPaymentMethodsCmd(), newSubscriptionsCmd())
	return cmd
}
