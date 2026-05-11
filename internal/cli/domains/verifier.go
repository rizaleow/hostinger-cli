package domains

import (
	"github.com/spf13/cobra"

	"github.com/rizaleow/hostinger-cli/internal/api"
	"github.com/rizaleow/hostinger-cli/internal/clictx"
)

func newVerifierCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "verifier", Short: "Domain Access Verifier"}
	cmd.AddCommand(newVerifierActiveCmd())
	return cmd
}

func newVerifierActiveCmd() *cobra.Command {
	var domains []string
	cmd := &cobra.Command{
		Use:   "active --domain <name>...",
		Short: "List active (pending/completed) verifications for given domains",
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := clictx.FromCommand(cmd).Client()
			if err != nil {
				return err
			}
			body := api.DomainAccessVerifierV2VerificationsListRequest{Domains: domains}
			resp, err := c.V2GetDomainVerificationsDIRECTWithResponse(cmd.Context(), body)
			if err != nil {
				return err
			}
			return clictx.Render(cmd, resp.JSON200)
		},
	}
	cmd.Flags().StringSliceVar(&domains, "domain", nil, "domain to check (repeatable)")
	_ = cmd.MarkFlagRequired("domain")
	return cmd
}
