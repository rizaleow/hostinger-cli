// Package dns implements `hostinger-cli dns` (zones + snapshots).
package dns

import "github.com/spf13/cobra"

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "dns", Short: "Manage DNS zones and snapshots"}
	cmd.AddCommand(newZoneCmd(), newSnapshotCmd())
	return cmd
}
