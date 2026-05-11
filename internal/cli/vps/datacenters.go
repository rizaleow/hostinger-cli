package vps

import (
	"github.com/spf13/cobra"

	"github.com/rizaleow/hostinger-cli/internal/clictx"
	"github.com/rizaleow/hostinger-cli/internal/output"
)

// dataCentersTableCols defines a friendly table view; falls back to JSON if a
// payload doesn't match.
var dataCentersTableCols = []output.Column{
	{Header: "ID", Path: "id"},
	{Header: "Name", Path: "name"},
	{Header: "City", Path: "city"},
	{Header: "Continent", Path: "continent"},
	{Header: "Location", Path: "location"},
}

func newDataCentersCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "data-centers", Short: "VPS data centers"}
	cmd.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List VPS data centers",
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := clictx.FromCommand(cmd).Client()
			if err != nil {
				return err
			}
			resp, err := c.VPSGetDataCenterListV1WithResponse(cmd.Context())
			if err != nil {
				return err
			}
			return clictx.Render(cmd, resp.JSON200)
		},
	})
	return cmd
}
