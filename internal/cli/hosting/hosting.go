// Package hosting implements `hostinger-cli hosting`.
package hosting

import (
	"github.com/spf13/cobra"

	"github.com/rizaleow/hostinger-cli/internal/api"
	"github.com/rizaleow/hostinger-cli/internal/clictx"
	"github.com/rizaleow/hostinger-cli/internal/cliutil"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "hosting", Short: "Hosting datacenters, domains, orders, websites"}
	cmd.AddCommand(newDatacentersCmd(), newDomainsCmd(), newOrdersCmd(), newWebsitesCmd())
	return cmd
}

func newDatacentersCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "datacenters", Short: "Hosting datacenters"}
	cmd.AddCommand(&cobra.Command{
		Use:   "list --order <id>",
		Short: "List available datacenters for an order",
		RunE: func(cmd *cobra.Command, _ []string) error {
			orderID, _ := cmd.Flags().GetInt("order")
			c, err := clictx.FromCommand(cmd).Client()
			if err != nil {
				return err
			}
			resp, err := c.HostingListAvailableDatacentersV1WithResponse(cmd.Context(),
				&api.HostingListAvailableDatacentersV1Params{OrderId: api.OrderIdRequired(orderID)})
			if err != nil {
				return err
			}
			return clictx.Render(cmd, resp.JSON200)
		},
	})
	cmd.Commands()[0].Flags().Int("order", 0, "order id (required)")
	_ = cmd.Commands()[0].MarkFlagRequired("order")
	return cmd
}

func newDomainsCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "domains", Short: "Hosting subdomains and ownership"}
	cmd.AddCommand(
		&cobra.Command{
			Use: "generate-subdomain", Short: "Generate a free Hostinger subdomain",
			RunE: func(cmd *cobra.Command, _ []string) error {
				c, err := clictx.FromCommand(cmd).Client()
				if err != nil {
					return err
				}
				resp, err := c.HostingGenerateAFreeSubdomainV1WithResponse(cmd.Context())
				if err != nil {
					return err
				}
				return clictx.Render(cmd, resp.JSON200)
			},
		},
		newVerifyOwnershipCmd(),
	)
	return cmd
}

func newVerifyOwnershipCmd() *cobra.Command {
	var domain string
	cmd := &cobra.Command{
		Use:   "verify-ownership --domain <name>",
		Short: "Verify ownership of a domain via TXT record",
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := clictx.FromCommand(cmd).Client()
			if err != nil {
				return err
			}
			body := api.HostingV1DomainsVerifyOwnershipRequest{Domain: domain}
			resp, err := c.HostingVerifyDomainOwnershipV1WithResponse(cmd.Context(), body)
			if err != nil {
				return err
			}
			return clictx.Render(cmd, resp.JSON200)
		},
	}
	cmd.Flags().StringVar(&domain, "domain", "", "domain to verify")
	_ = cmd.MarkFlagRequired("domain")
	return cmd
}

func newOrdersCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "orders", Short: "Hosting orders"}
	cmd.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List hosting orders",
		RunE: func(cmd *cobra.Command, _ []string) error {
			page, _ := cmd.Flags().GetInt("page")
			perPage, _ := cmd.Flags().GetInt("per-page")
			c, err := clictx.FromCommand(cmd).Client()
			if err != nil {
				return err
			}
			params := &api.HostingListOrdersV1Params{}
			if page > 0 {
				p := api.Page(page)
				params.Page = &p
			}
			if perPage > 0 {
				pp := api.PerPage(perPage)
				params.PerPage = &pp
			}
			resp, err := c.HostingListOrdersV1WithResponse(cmd.Context(), params)
			if err != nil {
				return err
			}
			return clictx.Render(cmd, resp.JSON200)
		},
	})
	cmd.Commands()[0].Flags().Int("page", 0, "page number")
	cmd.Commands()[0].Flags().Int("per-page", 0, "items per page")
	return cmd
}

func newWebsitesCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "websites", Short: "Hosted websites"}
	cmd.AddCommand(newWebsitesListCmd(), newWebsitesCreateCmd())
	return cmd
}

func newWebsitesListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List hosted websites",
		RunE: func(cmd *cobra.Command, _ []string) error {
			params := &api.HostingListWebsitesV1Params{}
			page, _ := cmd.Flags().GetInt("page")
			perPage, _ := cmd.Flags().GetInt("per-page")
			username, _ := cmd.Flags().GetString("username")
			orderID, _ := cmd.Flags().GetInt("order")
			enabled, _ := cmd.Flags().GetBool("enabled")
			enabledSet := cmd.Flags().Changed("enabled")
			domain, _ := cmd.Flags().GetString("domain")

			if page > 0 {
				p := api.Page(page)
				params.Page = &p
			}
			if perPage > 0 {
				pp := api.PerPage(perPage)
				params.PerPage = &pp
			}
			if username != "" {
				u := api.Username(username)
				params.Username = &u
			}
			if orderID > 0 {
				o := api.OrderId(orderID)
				params.OrderId = &o
			}
			if enabledSet {
				e := api.IsEnabled(enabled)
				params.IsEnabled = &e
			}
			if domain != "" {
				d := api.DomainFilter(domain)
				params.Domain = &d
			}
			c, err := clictx.FromCommand(cmd).Client()
			if err != nil {
				return err
			}
			resp, err := c.HostingListWebsitesV1WithResponse(cmd.Context(), params)
			if err != nil {
				return err
			}
			return clictx.Render(cmd, resp.JSON200)
		},
	}
	cmd.Flags().Int("page", 0, "page number")
	cmd.Flags().Int("per-page", 0, "items per page")
	cmd.Flags().String("username", "", "filter by hosting username")
	cmd.Flags().Int("order", 0, "filter by order id")
	cmd.Flags().Bool("enabled", false, "filter by enabled state")
	cmd.Flags().String("domain", "", "filter by domain name")
	return cmd
}

func newWebsitesCreateCmd() *cobra.Command {
	var fromFile string
	cmd := &cobra.Command{
		Use: "create --from-file <path>", Short: "Create a new hosted website",
		RunE: func(cmd *cobra.Command, _ []string) error {
			body, err := cliutil.ReadBody[api.HostingV1WebsitesCreateWebsiteRequest](fromFile)
			if err != nil {
				return err
			}
			c, err := clictx.FromCommand(cmd).Client()
			if err != nil {
				return err
			}
			resp, err := c.HostingCreateWebsiteV1WithResponse(cmd.Context(), body)
			if err != nil {
				return err
			}
			return clictx.Render(cmd, resp.JSON200)
		},
	}
	cmd.Flags().StringVar(&fromFile, "from-file", "", "JSON body file ('-' for stdin)")
	_ = cmd.MarkFlagRequired("from-file")
	return cmd
}
