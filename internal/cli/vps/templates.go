package vps

import (
	"github.com/spf13/cobra"

	"github.com/rizaleow/hostinger-cli/internal/api"
	"github.com/rizaleow/hostinger-cli/internal/clictx"
	"github.com/rizaleow/hostinger-cli/internal/cliutil"
)

func newTemplatesCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "templates", Short: "VPS OS templates"}
	cmd.AddCommand(
		&cobra.Command{
			Use: "list", Short: "List templates",
			RunE: func(cmd *cobra.Command, _ []string) error {
				c, err := clictx.FromCommand(cmd).Client()
				if err != nil {
					return err
				}
				resp, err := c.VPSGetTemplatesV1WithResponse(cmd.Context())
				if err != nil {
					return err
				}
				return clictx.Render(cmd, resp.JSON200)
			},
		},
		&cobra.Command{
			Use: "get <id>", Short: "Get a template", Args: cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				id, err := cliutil.ParseInt(args[0])
				if err != nil {
					return err
				}
				c, err := clictx.FromCommand(cmd).Client()
				if err != nil {
					return err
				}
				resp, err := c.VPSGetTemplateDetailsV1WithResponse(cmd.Context(), api.TemplateId(id))
				if err != nil {
					return err
				}
				return clictx.Render(cmd, resp.JSON200)
			},
		},
	)
	return cmd
}
