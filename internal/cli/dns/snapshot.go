package dns

import (
	"strconv"

	"github.com/spf13/cobra"

	"github.com/rizaleow/hostinger-cli/internal/api"
	"github.com/rizaleow/hostinger-cli/internal/clictx"
)

func newSnapshotCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "snapshot", Short: "DNS snapshots"}
	cmd.AddCommand(newSnapshotListCmd(), newSnapshotGetCmd(), newSnapshotRestoreCmd())
	return cmd
}

func newSnapshotListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list <domain>",
		Short: "List DNS snapshots for a domain",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clictx.FromCommand(cmd).Client()
			if err != nil {
				return err
			}
			resp, err := client.DNSGetDNSSnapshotListV1WithResponse(cmd.Context(), api.Domain(args[0]))
			if err != nil {
				return err
			}
			return clictx.Render(cmd, resp.JSON200)
		},
	}
}

func newSnapshotGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <domain> <snapshot-id>",
		Short: "Get a specific DNS snapshot",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.Atoi(args[1])
			if err != nil {
				return err
			}
			client, err := clictx.FromCommand(cmd).Client()
			if err != nil {
				return err
			}
			resp, err := client.DNSGetDNSSnapshotV1WithResponse(cmd.Context(), api.Domain(args[0]), api.SnapshotId(id))
			if err != nil {
				return err
			}
			return clictx.Render(cmd, resp.JSON200)
		},
	}
}

func newSnapshotRestoreCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "restore <domain> <snapshot-id>",
		Short: "Restore a DNS snapshot",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.Atoi(args[1])
			if err != nil {
				return err
			}
			client, err := clictx.FromCommand(cmd).Client()
			if err != nil {
				return err
			}
			resp, err := client.DNSRestoreDNSSnapshotV1WithResponse(cmd.Context(), api.Domain(args[0]), api.SnapshotId(id))
			if err != nil {
				return err
			}
			return clictx.Render(cmd, resp.JSON200)
		},
	}
}
