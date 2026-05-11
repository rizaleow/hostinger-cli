package domains

import (
	"github.com/spf13/cobra"

	"github.com/rizaleow/hostinger-cli/internal/api"
	"github.com/rizaleow/hostinger-cli/internal/clictx"
	"github.com/rizaleow/hostinger-cli/internal/cliutil"
)

func newForwardingCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "forwarding", Short: "Domain forwarding"}
	cmd.AddCommand(newForwardingGetCmd(), newForwardingCreateCmd(), newForwardingDeleteCmd())
	return cmd
}

func newForwardingGetCmd() *cobra.Command {
	return &cobra.Command{
		Use: "get <domain>", Short: "Get forwarding config", Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := clictx.FromCommand(cmd).Client()
			if err != nil {
				return err
			}
			resp, err := c.DomainsGetDomainForwardingV1WithResponse(cmd.Context(), api.Domain(args[0]))
			if err != nil {
				return err
			}
			return clictx.Render(cmd, resp.JSON200)
		},
	}
}

func newForwardingCreateCmd() *cobra.Command {
	var fromFile string
	cmd := &cobra.Command{
		Use: "create --from-file <path>", Short: "Create forwarding config",
		RunE: func(cmd *cobra.Command, _ []string) error {
			body, err := cliutil.ReadBody[api.DomainsV1ForwardingStoreRequest](fromFile)
			if err != nil {
				return err
			}
			c, err := clictx.FromCommand(cmd).Client()
			if err != nil {
				return err
			}
			resp, err := c.DomainsCreateDomainForwardingV1WithResponse(cmd.Context(), body)
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

func newForwardingDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use: "delete <domain>", Short: "Remove forwarding config", Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := clictx.FromCommand(cmd).Client()
			if err != nil {
				return err
			}
			resp, err := c.DomainsDeleteDomainForwardingV1WithResponse(cmd.Context(), api.Domain(args[0]))
			if err != nil {
				return err
			}
			return clictx.Render(cmd, resp.JSON200)
		},
	}
}
