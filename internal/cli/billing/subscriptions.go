package billing

import (
	"github.com/spf13/cobra"

	"github.com/rizaleow/hostinger-cli/internal/api"
	"github.com/rizaleow/hostinger-cli/internal/clictx"
)

func newSubscriptionsCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "subscriptions", Short: "Manage subscriptions"}
	cmd.AddCommand(
		newSubscriptionsListCmd(),
		newSubscriptionsEnableAutoRenewalCmd(),
		newSubscriptionsDisableAutoRenewalCmd(),
	)
	return cmd
}

func newSubscriptionsListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List subscriptions",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clictx.FromCommand(cmd).Client()
			if err != nil {
				return err
			}
			resp, err := client.BillingGetSubscriptionListV1WithResponse(cmd.Context())
			if err != nil {
				return err
			}
			return clictx.Render(cmd, resp.JSON200)
		},
	}
}

func newSubscriptionsEnableAutoRenewalCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "enable-auto-renewal <subscription-id>",
		Short: "Enable auto-renewal for a subscription",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clictx.FromCommand(cmd).Client()
			if err != nil {
				return err
			}
			resp, err := client.BillingEnableAutoRenewalV1WithResponse(cmd.Context(), api.SubscriptionId(args[0]))
			if err != nil {
				return err
			}
			return clictx.Render(cmd, resp.JSON200)
		},
	}
}

func newSubscriptionsDisableAutoRenewalCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "disable-auto-renewal <subscription-id>",
		Short: "Disable auto-renewal for a subscription",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clictx.FromCommand(cmd).Client()
			if err != nil {
				return err
			}
			resp, err := client.BillingDisableAutoRenewalV1WithResponse(cmd.Context(), api.SubscriptionId(args[0]))
			if err != nil {
				return err
			}
			return clictx.Render(cmd, resp.JSON200)
		},
	}
}
