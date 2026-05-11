package billing

import (
	"github.com/spf13/cobra"

	"github.com/rizaleow/hostinger-cli/internal/api"
	"github.com/rizaleow/hostinger-cli/internal/clictx"
)

func newCatalogCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "catalog", Short: "Catalog items available for order"}
	cmd.AddCommand(newCatalogListCmd())
	return cmd
}

func newCatalogListCmd() *cobra.Command {
	var (
		category string
		name     string
	)
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List catalog items",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clictx.FromCommand(cmd).Client()
			if err != nil {
				return err
			}
			params := &api.BillingGetCatalogItemListV1Params{}
			if category != "" {
				c := api.BillingGetCatalogItemListV1ParamsCategory(category)
				params.Category = &c
			}
			if name != "" {
				n := api.Name(name)
				params.Name = &n
			}
			resp, err := client.BillingGetCatalogItemListV1WithResponse(cmd.Context(), params)
			if err != nil {
				return err
			}
			return clictx.Render(cmd, resp.JSON200)
		},
	}
	cmd.Flags().StringVar(&category, "category", "", "filter by category")
	cmd.Flags().StringVar(&name, "name", "", "filter by name")
	return cmd
}
