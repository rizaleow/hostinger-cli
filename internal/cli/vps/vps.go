// Package vps implements `hostinger-cli vps`.
package vps

import "github.com/spf13/cobra"

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "vps", Short: "VPS instances, networking, snapshots, and Docker"}
	cmd.AddCommand(
		newDataCentersCmd(),
		newVMCmd(),
		newDockerCmd(),
		newFirewallCmd(),
		newPostInstallScriptsCmd(),
		newPublicKeysCmd(),
		newTemplatesCmd(),
	)
	return cmd
}
