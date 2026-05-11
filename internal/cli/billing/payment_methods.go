package billing

import (
	"github.com/spf13/cobra"

	"github.com/rizaleow/hostinger-cli/internal/api"
	"github.com/rizaleow/hostinger-cli/internal/clictx"
)

func newPaymentMethodsCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "payment-methods", Short: "Manage payment methods"}
	cmd.AddCommand(
		newPaymentMethodsListCmd(),
		newPaymentMethodsSetDefaultCmd(),
		newPaymentMethodsDeleteCmd(),
	)
	return cmd
}

func newPaymentMethodsListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List payment methods",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clictx.FromCommand(cmd).Client()
			if err != nil {
				return err
			}
			resp, err := client.BillingGetPaymentMethodListV1WithResponse(cmd.Context())
			if err != nil {
				return err
			}
			return clictx.Render(cmd, resp.JSON200)
		},
	}
}

func newPaymentMethodsSetDefaultCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "set-default <payment-method-id>",
		Short: "Make a payment method the default",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clictx.FromCommand(cmd).Client()
			if err != nil {
				return err
			}
			id, err := parseInt64(args[0])
			if err != nil {
				return err
			}
			resp, err := client.BillingSetDefaultPaymentMethodV1WithResponse(cmd.Context(), api.PaymentMethodId(id))
			if err != nil {
				return err
			}
			return clictx.Render(cmd, resp.JSON200)
		},
	}
}

func newPaymentMethodsDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <payment-method-id>",
		Short: "Delete a payment method",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clictx.FromCommand(cmd).Client()
			if err != nil {
				return err
			}
			id, err := parseInt64(args[0])
			if err != nil {
				return err
			}
			resp, err := client.BillingDeletePaymentMethodV1WithResponse(cmd.Context(), api.PaymentMethodId(id))
			if err != nil {
				return err
			}
			return clictx.Render(cmd, resp.JSON200)
		},
	}
}
