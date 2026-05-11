package domains

import (
	"github.com/spf13/cobra"

	"github.com/rizaleow/hostinger-cli/internal/api"
	"github.com/rizaleow/hostinger-cli/internal/clictx"
)

func newAvailabilityCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "availability", Short: "Domain availability"}
	cmd.AddCommand(newAvailabilityCheckCmd())
	return cmd
}

func newAvailabilityCheckCmd() *cobra.Command {
	var (
		domain           string
		tlds             []string
		withAlternatives bool
	)
	cmd := &cobra.Command{
		Use:   "check --domain <name> --tld <tld>...",
		Short: "Check availability of a domain across TLDs",
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := clictx.FromCommand(cmd).Client()
			if err != nil {
				return err
			}
			body := api.DomainsV1AvailabilityAvailabilityRequest{
				Domain:           domain,
				Tlds:             tlds,
				WithAlternatives: &withAlternatives,
			}
			resp, err := c.DomainsCheckDomainAvailabilityV1WithResponse(cmd.Context(), body)
			if err != nil {
				return err
			}
			return clictx.Render(cmd, resp.JSON200)
		},
	}
	cmd.Flags().StringVar(&domain, "domain", "", "domain name without TLD (e.g. mysite)")
	cmd.Flags().StringSliceVar(&tlds, "tld", nil, "TLD without leading dot (repeatable)")
	cmd.Flags().BoolVar(&withAlternatives, "with-alternatives", false, "include alternative suggestions (single TLD only)")
	_ = cmd.MarkFlagRequired("domain")
	_ = cmd.MarkFlagRequired("tld")
	return cmd
}
