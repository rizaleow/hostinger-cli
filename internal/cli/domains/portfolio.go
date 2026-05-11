package domains

import (
	"github.com/spf13/cobra"

	"github.com/rizaleow/hostinger-cli/internal/api"
	"github.com/rizaleow/hostinger-cli/internal/clictx"
	"github.com/rizaleow/hostinger-cli/internal/cliutil"
)

func newPortfolioCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "portfolio", Short: "Manage owned domains"}
	cmd.AddCommand(
		newPortfolioListCmd(),
		newPortfolioGetCmd(),
		newPortfolioPurchaseCmd(),
		newPortfolioLockCmd(),
		newPortfolioUnlockCmd(),
		newPortfolioPrivacyEnableCmd(),
		newPortfolioPrivacyDisableCmd(),
		newPortfolioNameserversCmd(),
	)
	return cmd
}

func newPortfolioListCmd() *cobra.Command {
	return &cobra.Command{
		Use: "list", Short: "List domains",
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := clictx.FromCommand(cmd).Client()
			if err != nil {
				return err
			}
			resp, err := c.DomainsGetDomainListV1WithResponse(cmd.Context())
			if err != nil {
				return err
			}
			return clictx.Render(cmd, resp.JSON200)
		},
	}
}

func newPortfolioGetCmd() *cobra.Command {
	return &cobra.Command{
		Use: "get <domain>", Short: "Get domain details", Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := clictx.FromCommand(cmd).Client()
			if err != nil {
				return err
			}
			resp, err := c.DomainsGetDomainDetailsV1WithResponse(cmd.Context(), api.Domain(args[0]))
			if err != nil {
				return err
			}
			return clictx.Render(cmd, resp.JSON200)
		},
	}
}

func newPortfolioPurchaseCmd() *cobra.Command {
	var fromFile string
	cmd := &cobra.Command{
		Use: "purchase --from-file <path>", Short: "Purchase a new domain",
		RunE: func(cmd *cobra.Command, _ []string) error {
			body, err := cliutil.ReadBody[api.DomainsV1PortfolioPurchaseRequest](fromFile)
			if err != nil {
				return err
			}
			c, err := clictx.FromCommand(cmd).Client()
			if err != nil {
				return err
			}
			resp, err := c.DomainsPurchaseNewDomainV1WithResponse(cmd.Context(), body)
			if err != nil {
				return err
			}
			return clictx.Render(cmd, resp.JSON200)
		},
	}
	cmd.Flags().StringVar(&fromFile, "from-file", "", "JSON request body file ('-' for stdin)")
	_ = cmd.MarkFlagRequired("from-file")
	return cmd
}

func newPortfolioLockCmd() *cobra.Command {
	return &cobra.Command{
		Use: "lock <domain>", Short: "Enable domain transfer lock", Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := clictx.FromCommand(cmd).Client()
			if err != nil {
				return err
			}
			resp, err := c.DomainsEnableDomainLockV1WithResponse(cmd.Context(), api.Domain(args[0]))
			if err != nil {
				return err
			}
			return clictx.Render(cmd, resp.JSON200)
		},
	}
}

func newPortfolioUnlockCmd() *cobra.Command {
	return &cobra.Command{
		Use: "unlock <domain>", Short: "Disable domain transfer lock", Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := clictx.FromCommand(cmd).Client()
			if err != nil {
				return err
			}
			resp, err := c.DomainsDisableDomainLockV1WithResponse(cmd.Context(), api.Domain(args[0]))
			if err != nil {
				return err
			}
			return clictx.Render(cmd, resp.JSON200)
		},
	}
}

func newPortfolioPrivacyEnableCmd() *cobra.Command {
	return &cobra.Command{
		Use: "enable-privacy <domain>", Short: "Enable WHOIS privacy protection", Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := clictx.FromCommand(cmd).Client()
			if err != nil {
				return err
			}
			resp, err := c.DomainsEnablePrivacyProtectionV1WithResponse(cmd.Context(), api.Domain(args[0]))
			if err != nil {
				return err
			}
			return clictx.Render(cmd, resp.JSON200)
		},
	}
}

func newPortfolioPrivacyDisableCmd() *cobra.Command {
	return &cobra.Command{
		Use: "disable-privacy <domain>", Short: "Disable WHOIS privacy protection", Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := clictx.FromCommand(cmd).Client()
			if err != nil {
				return err
			}
			resp, err := c.DomainsDisablePrivacyProtectionV1WithResponse(cmd.Context(), api.Domain(args[0]))
			if err != nil {
				return err
			}
			return clictx.Render(cmd, resp.JSON200)
		},
	}
}

func newPortfolioNameserversCmd() *cobra.Command {
	var fromFile string
	cmd := &cobra.Command{
		Use: "update-nameservers <domain> --from-file <path>", Short: "Update nameservers", Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			body, err := cliutil.ReadBody[api.DomainsV1PortfolioUpdateNameserversRequest](fromFile)
			if err != nil {
				return err
			}
			c, err := clictx.FromCommand(cmd).Client()
			if err != nil {
				return err
			}
			resp, err := c.DomainsUpdateDomainNameserversV1WithResponse(cmd.Context(), api.Domain(args[0]), body)
			if err != nil {
				return err
			}
			return clictx.Render(cmd, resp.JSON200)
		},
	}
	cmd.Flags().StringVar(&fromFile, "from-file", "", "JSON request body file ('-' for stdin)")
	_ = cmd.MarkFlagRequired("from-file")
	return cmd
}
