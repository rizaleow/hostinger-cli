package domains

import (
	"github.com/spf13/cobra"

	"github.com/rizaleow/hostinger-cli/internal/api"
	"github.com/rizaleow/hostinger-cli/internal/clictx"
	"github.com/rizaleow/hostinger-cli/internal/cliutil"
)

func newWHOISCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "whois", Short: "WHOIS contact profiles"}
	cmd.AddCommand(
		newWHOISListCmd(),
		newWHOISGetCmd(),
		newWHOISCreateCmd(),
		newWHOISDeleteCmd(),
		newWHOISUsageCmd(),
	)
	return cmd
}

func newWHOISListCmd() *cobra.Command {
	var tld string
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List WHOIS profiles",
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := clictx.FromCommand(cmd).Client()
			if err != nil {
				return err
			}
			params := &api.DomainsGetWHOISProfileListV1Params{}
			if tld != "" {
				t := api.Tld(tld)
				params.Tld = &t
			}
			resp, err := c.DomainsGetWHOISProfileListV1WithResponse(cmd.Context(), params)
			if err != nil {
				return err
			}
			return clictx.Render(cmd, resp.JSON200)
		},
	}
	cmd.Flags().StringVar(&tld, "tld", "", "filter by TLD (without leading dot)")
	return cmd
}

func newWHOISGetCmd() *cobra.Command {
	return &cobra.Command{
		Use: "get <whois-id>", Short: "Get WHOIS profile", Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := cliutil.ParseInt(args[0])
			if err != nil {
				return err
			}
			c, err := clictx.FromCommand(cmd).Client()
			if err != nil {
				return err
			}
			resp, err := c.DomainsGetWHOISProfileV1WithResponse(cmd.Context(), api.WhoisId(id))
			if err != nil {
				return err
			}
			return clictx.Render(cmd, resp.JSON200)
		},
	}
}

func newWHOISCreateCmd() *cobra.Command {
	var fromFile string
	cmd := &cobra.Command{
		Use: "create --from-file <path>", Short: "Create a WHOIS profile",
		RunE: func(cmd *cobra.Command, _ []string) error {
			body, err := cliutil.ReadBody[api.DomainsV1WHOISStoreRequest](fromFile)
			if err != nil {
				return err
			}
			c, err := clictx.FromCommand(cmd).Client()
			if err != nil {
				return err
			}
			resp, err := c.DomainsCreateWHOISProfileV1WithResponse(cmd.Context(), body)
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

func newWHOISDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use: "delete <whois-id>", Short: "Delete a WHOIS profile", Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := cliutil.ParseInt(args[0])
			if err != nil {
				return err
			}
			c, err := clictx.FromCommand(cmd).Client()
			if err != nil {
				return err
			}
			resp, err := c.DomainsDeleteWHOISProfileV1WithResponse(cmd.Context(), api.WhoisId(id))
			if err != nil {
				return err
			}
			return clictx.Render(cmd, resp.JSON200)
		},
	}
}

func newWHOISUsageCmd() *cobra.Command {
	return &cobra.Command{
		Use: "usage <whois-id>", Short: "List domains using a WHOIS profile", Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := cliutil.ParseInt(args[0])
			if err != nil {
				return err
			}
			c, err := clictx.FromCommand(cmd).Client()
			if err != nil {
				return err
			}
			resp, err := c.DomainsGetWHOISProfileUsageV1WithResponse(cmd.Context(), api.WhoisId(id))
			if err != nil {
				return err
			}
			return clictx.Render(cmd, resp.JSON200)
		},
	}
}
