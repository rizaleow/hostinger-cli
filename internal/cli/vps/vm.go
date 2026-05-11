package vps

import (
	"time"

	"github.com/spf13/cobra"

	"github.com/rizaleow/hostinger-cli/internal/api"
	"github.com/rizaleow/hostinger-cli/internal/clictx"
	"github.com/rizaleow/hostinger-cli/internal/cliutil"
)

func newVMCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "vm", Short: "Virtual machine instances"}
	cmd.AddCommand(
		newVMListCmd(),
		newVMGetCmd(),
		newVMPurchaseCmd(),
		newVMSetupCmd(),
		newVMStartCmd(),
		newVMStopCmd(),
		newVMRestartCmd(),
		newVMRecreateCmd(),
		newVMSetHostnameCmd(),
		newVMResetHostnameCmd(),
		newVMSetRootPasswordCmd(),
		newVMSetPanelPasswordCmd(),
		newVMSetNameserversCmd(),
		newVMMetricsCmd(),
		newVMActionsCmd(),
		newVMActionGetCmd(),
		newVMAttachedKeysCmd(),
		newVMBackupsCmd(),
		newVMBackupRestoreCmd(),
		newVMSnapshotCmd(),
		newVMSnapshotGetCmd(),
		newVMSnapshotRestoreCmd(),
		newVMSnapshotDeleteCmd(),
		newVMRecoveryStartCmd(),
		newVMRecoveryStopCmd(),
		newVMPTRCreateCmd(),
		newVMPTRDeleteCmd(),
		newVMMonarxInstallCmd(),
		newVMMonarxUninstallCmd(),
		newVMMonarxMetricsCmd(),
	)
	return cmd
}

func vmID(arg string) (api.VirtualMachineId, error) {
	n, err := cliutil.ParseInt(arg)
	if err != nil {
		return 0, err
	}
	return api.VirtualMachineId(n), nil
}

func vmCmdNoBody(use, short string, do func(c *api.ClientWithResponses, ctx, vmID any) (any, error)) *cobra.Command {
	// helper unused — kept simple inline implementations below.
	return nil
}

// --- read-only ---

func newVMListCmd() *cobra.Command {
	return &cobra.Command{
		Use: "list", Short: "List VPS virtual machines",
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := clictx.FromCommand(cmd).Client()
			if err != nil {
				return err
			}
			resp, err := c.VPSGetVirtualMachinesV1WithResponse(cmd.Context())
			if err != nil {
				return err
			}
			return clictx.Render(cmd, resp.JSON200)
		},
	}
}

func newVMGetCmd() *cobra.Command {
	return &cobra.Command{
		Use: "get <vm-id>", Short: "Get VM details", Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := vmID(args[0])
			if err != nil {
				return err
			}
			c, err := clictx.FromCommand(cmd).Client()
			if err != nil {
				return err
			}
			resp, err := c.VPSGetVirtualMachineDetailsV1WithResponse(cmd.Context(), id)
			if err != nil {
				return err
			}
			return clictx.Render(cmd, resp.JSON200)
		},
	}
}

func newVMMetricsCmd() *cobra.Command {
	var from, to string
	cmd := &cobra.Command{
		Use: "metrics <vm-id>", Short: "Get VM metrics (RFC3339 date window)", Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := vmID(args[0])
			if err != nil {
				return err
			}
			fromT, err := time.Parse(time.RFC3339, from)
			if err != nil {
				return err
			}
			toT, err := time.Parse(time.RFC3339, to)
			if err != nil {
				return err
			}
			c, err := clictx.FromCommand(cmd).Client()
			if err != nil {
				return err
			}
			params := &api.VPSGetMetricsV1Params{DateFrom: fromT, DateTo: toT}
			resp, err := c.VPSGetMetricsV1WithResponse(cmd.Context(), id, params)
			if err != nil {
				return err
			}
			return clictx.Render(cmd, resp.JSON200)
		},
	}
	cmd.Flags().StringVar(&from, "from", "", "RFC3339 start (e.g. 2026-01-01T00:00:00Z)")
	cmd.Flags().StringVar(&to, "to", "", "RFC3339 end")
	_ = cmd.MarkFlagRequired("from")
	_ = cmd.MarkFlagRequired("to")
	return cmd
}

func newVMActionsCmd() *cobra.Command {
	return &cobra.Command{
		Use: "actions <vm-id>", Short: "List VM actions", Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := vmID(args[0])
			if err != nil {
				return err
			}
			c, err := clictx.FromCommand(cmd).Client()
			if err != nil {
				return err
			}
			resp, err := c.VPSGetActionsV1WithResponse(cmd.Context(), id, &api.VPSGetActionsV1Params{})
			if err != nil {
				return err
			}
			return clictx.Render(cmd, resp.JSON200)
		},
	}
}

func newVMActionGetCmd() *cobra.Command {
	return &cobra.Command{
		Use: "action-get <vm-id> <action-id>", Short: "Get a single VM action", Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := vmID(args[0])
			if err != nil {
				return err
			}
			aid, err := cliutil.ParseInt(args[1])
			if err != nil {
				return err
			}
			c, err := clictx.FromCommand(cmd).Client()
			if err != nil {
				return err
			}
			resp, err := c.VPSGetActionDetailsV1WithResponse(cmd.Context(), id, api.ActionId(aid))
			if err != nil {
				return err
			}
			return clictx.Render(cmd, resp.JSON200)
		},
	}
}

func newVMAttachedKeysCmd() *cobra.Command {
	return &cobra.Command{
		Use: "attached-keys <vm-id>", Short: "List public keys attached to a VM", Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := vmID(args[0])
			if err != nil {
				return err
			}
			c, err := clictx.FromCommand(cmd).Client()
			if err != nil {
				return err
			}
			resp, err := c.VPSGetAttachedPublicKeysV1WithResponse(cmd.Context(), id, &api.VPSGetAttachedPublicKeysV1Params{})
			if err != nil {
				return err
			}
			return clictx.Render(cmd, resp.JSON200)
		},
	}
}

func newVMBackupsCmd() *cobra.Command {
	return &cobra.Command{
		Use: "backups <vm-id>", Short: "List VM backups", Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := vmID(args[0])
			if err != nil {
				return err
			}
			c, err := clictx.FromCommand(cmd).Client()
			if err != nil {
				return err
			}
			resp, err := c.VPSGetBackupsV1WithResponse(cmd.Context(), id, &api.VPSGetBackupsV1Params{})
			if err != nil {
				return err
			}
			return clictx.Render(cmd, resp.JSON200)
		},
	}
}

func newVMBackupRestoreCmd() *cobra.Command {
	return &cobra.Command{
		Use: "backup-restore <vm-id> <backup-id>", Short: "Restore a VM backup", Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := vmID(args[0])
			if err != nil {
				return err
			}
			b, err := cliutil.ParseInt(args[1])
			if err != nil {
				return err
			}
			c, err := clictx.FromCommand(cmd).Client()
			if err != nil {
				return err
			}
			resp, err := c.VPSRestoreBackupV1WithResponse(cmd.Context(), id, api.BackupId(b))
			if err != nil {
				return err
			}
			return clictx.Render(cmd, resp.JSON200)
		},
	}
}

// --- mutations ---

func newVMPurchaseCmd() *cobra.Command {
	var fromFile string
	cmd := &cobra.Command{
		Use: "purchase --from-file <path>", Short: "Purchase a new VM",
		RunE: func(cmd *cobra.Command, _ []string) error {
			body, err := cliutil.ReadBody[api.VPSV1VirtualMachinePurchaseRequest](fromFile)
			if err != nil {
				return err
			}
			c, err := clictx.FromCommand(cmd).Client()
			if err != nil {
				return err
			}
			resp, err := c.VPSPurchaseNewVirtualMachineV1WithResponse(cmd.Context(), body)
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

func newVMSetupCmd() *cobra.Command {
	var fromFile string
	cmd := &cobra.Command{
		Use: "setup <vm-id> --from-file <path>", Short: "Configure a purchased VM", Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := vmID(args[0])
			if err != nil {
				return err
			}
			body, err := cliutil.ReadBody[api.VPSV1VirtualMachineSetupRequest](fromFile)
			if err != nil {
				return err
			}
			c, err := clictx.FromCommand(cmd).Client()
			if err != nil {
				return err
			}
			resp, err := c.VPSSetupPurchasedVirtualMachineV1WithResponse(cmd.Context(), id, body)
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

func newVMStartCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "start <vm-id>", Short: "Start a VM", Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := vmID(args[0])
			if err != nil {
				return err
			}
			c, err := clictx.FromCommand(cmd).Client()
			if err != nil {
				return err
			}
			resp, err := c.VPSStartVirtualMachineV1WithResponse(cmd.Context(), id)
			if err != nil {
				return err
			}
			return maybeWait(cmd, id, resp.JSON200)
		},
	}
	addWaitFlags(cmd)
	return cmd
}

func newVMStopCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "stop <vm-id>", Short: "Stop a VM", Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := vmID(args[0])
			if err != nil {
				return err
			}
			c, err := clictx.FromCommand(cmd).Client()
			if err != nil {
				return err
			}
			resp, err := c.VPSStopVirtualMachineV1WithResponse(cmd.Context(), id)
			if err != nil {
				return err
			}
			return maybeWait(cmd, id, resp.JSON200)
		},
	}
	addWaitFlags(cmd)
	return cmd
}

func newVMRestartCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "restart <vm-id>", Short: "Restart a VM", Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := vmID(args[0])
			if err != nil {
				return err
			}
			c, err := clictx.FromCommand(cmd).Client()
			if err != nil {
				return err
			}
			resp, err := c.VPSRestartVirtualMachineV1WithResponse(cmd.Context(), id)
			if err != nil {
				return err
			}
			return maybeWait(cmd, id, resp.JSON200)
		},
	}
	addWaitFlags(cmd)
	return cmd
}

func newVMRecreateCmd() *cobra.Command {
	var fromFile string
	cmd := &cobra.Command{
		Use: "recreate <vm-id> --from-file <path>", Short: "Recreate a VM with a new image",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := vmID(args[0])
			if err != nil {
				return err
			}
			body, err := cliutil.ReadBody[api.VPSV1VirtualMachineRecreateRequest](fromFile)
			if err != nil {
				return err
			}
			c, err := clictx.FromCommand(cmd).Client()
			if err != nil {
				return err
			}
			resp, err := c.VPSRecreateVirtualMachineV1WithResponse(cmd.Context(), id, body)
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

func newVMSetHostnameCmd() *cobra.Command {
	var fromFile string
	cmd := &cobra.Command{
		Use: "set-hostname <vm-id> --from-file <path>", Short: "Set VM hostname", Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := vmID(args[0])
			if err != nil {
				return err
			}
			body, err := cliutil.ReadBody[api.VPSV1VirtualMachineHostnameUpdateRequest](fromFile)
			if err != nil {
				return err
			}
			c, err := clictx.FromCommand(cmd).Client()
			if err != nil {
				return err
			}
			resp, err := c.VPSSetHostnameV1WithResponse(cmd.Context(), id, body)
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

func newVMResetHostnameCmd() *cobra.Command {
	return &cobra.Command{
		Use: "reset-hostname <vm-id>", Short: "Reset VM hostname to default", Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := vmID(args[0])
			if err != nil {
				return err
			}
			c, err := clictx.FromCommand(cmd).Client()
			if err != nil {
				return err
			}
			resp, err := c.VPSResetHostnameV1WithResponse(cmd.Context(), id)
			if err != nil {
				return err
			}
			return clictx.Render(cmd, resp.JSON200)
		},
	}
}

func newVMSetRootPasswordCmd() *cobra.Command {
	var fromFile string
	cmd := &cobra.Command{
		Use: "set-root-password <vm-id> --from-file <path>", Short: "Set VM root password",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := vmID(args[0])
			if err != nil {
				return err
			}
			body, err := cliutil.ReadBody[api.VPSV1VirtualMachineRootPasswordUpdateRequest](fromFile)
			if err != nil {
				return err
			}
			c, err := clictx.FromCommand(cmd).Client()
			if err != nil {
				return err
			}
			resp, err := c.VPSSetRootPasswordV1WithResponse(cmd.Context(), id, body)
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

func newVMSetPanelPasswordCmd() *cobra.Command {
	var fromFile string
	cmd := &cobra.Command{
		Use: "set-panel-password <vm-id> --from-file <path>", Short: "Set VM panel password",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := vmID(args[0])
			if err != nil {
				return err
			}
			body, err := cliutil.ReadBody[api.VPSV1VirtualMachinePanelPasswordUpdateRequest](fromFile)
			if err != nil {
				return err
			}
			c, err := clictx.FromCommand(cmd).Client()
			if err != nil {
				return err
			}
			resp, err := c.VPSSetPanelPasswordV1WithResponse(cmd.Context(), id, body)
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

func newVMSetNameserversCmd() *cobra.Command {
	var fromFile string
	cmd := &cobra.Command{
		Use: "set-nameservers <vm-id> --from-file <path>", Short: "Set VM nameservers",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := vmID(args[0])
			if err != nil {
				return err
			}
			body, err := cliutil.ReadBody[api.VPSV1VirtualMachineNameserversUpdateRequest](fromFile)
			if err != nil {
				return err
			}
			c, err := clictx.FromCommand(cmd).Client()
			if err != nil {
				return err
			}
			resp, err := c.VPSSetNameserversV1WithResponse(cmd.Context(), id, body)
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

func newVMSnapshotCmd() *cobra.Command {
	return &cobra.Command{
		Use: "snapshot <vm-id>", Short: "Create a VM snapshot", Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := vmID(args[0])
			if err != nil {
				return err
			}
			c, err := clictx.FromCommand(cmd).Client()
			if err != nil {
				return err
			}
			resp, err := c.VPSCreateSnapshotV1WithResponse(cmd.Context(), id)
			if err != nil {
				return err
			}
			return clictx.Render(cmd, resp.JSON200)
		},
	}
}

func newVMSnapshotGetCmd() *cobra.Command {
	return &cobra.Command{
		Use: "snapshot-get <vm-id>", Short: "Get the current VM snapshot", Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := vmID(args[0])
			if err != nil {
				return err
			}
			c, err := clictx.FromCommand(cmd).Client()
			if err != nil {
				return err
			}
			resp, err := c.VPSGetSnapshotV1WithResponse(cmd.Context(), id)
			if err != nil {
				return err
			}
			return clictx.Render(cmd, resp.JSON200)
		},
	}
}

func newVMSnapshotRestoreCmd() *cobra.Command {
	return &cobra.Command{
		Use: "snapshot-restore <vm-id>", Short: "Restore the current VM snapshot",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := vmID(args[0])
			if err != nil {
				return err
			}
			c, err := clictx.FromCommand(cmd).Client()
			if err != nil {
				return err
			}
			resp, err := c.VPSRestoreSnapshotV1WithResponse(cmd.Context(), id)
			if err != nil {
				return err
			}
			return clictx.Render(cmd, resp.JSON200)
		},
	}
}

func newVMSnapshotDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use: "snapshot-delete <vm-id>", Short: "Delete the current VM snapshot",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := vmID(args[0])
			if err != nil {
				return err
			}
			c, err := clictx.FromCommand(cmd).Client()
			if err != nil {
				return err
			}
			resp, err := c.VPSDeleteSnapshotV1WithResponse(cmd.Context(), id)
			if err != nil {
				return err
			}
			return clictx.Render(cmd, resp.JSON200)
		},
	}
}

func newVMRecoveryStartCmd() *cobra.Command {
	var fromFile string
	cmd := &cobra.Command{
		Use: "recovery-start <vm-id> --from-file <path>", Short: "Start VM recovery mode",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := vmID(args[0])
			if err != nil {
				return err
			}
			body, err := cliutil.ReadBody[api.VPSV1VirtualMachineRecoveryStartRequest](fromFile)
			if err != nil {
				return err
			}
			c, err := clictx.FromCommand(cmd).Client()
			if err != nil {
				return err
			}
			resp, err := c.VPSStartRecoveryModeV1WithResponse(cmd.Context(), id, body)
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

func newVMRecoveryStopCmd() *cobra.Command {
	return &cobra.Command{
		Use: "recovery-stop <vm-id>", Short: "Stop VM recovery mode", Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := vmID(args[0])
			if err != nil {
				return err
			}
			c, err := clictx.FromCommand(cmd).Client()
			if err != nil {
				return err
			}
			resp, err := c.VPSStopRecoveryModeV1WithResponse(cmd.Context(), id)
			if err != nil {
				return err
			}
			return clictx.Render(cmd, resp.JSON200)
		},
	}
}

func newVMPTRCreateCmd() *cobra.Command {
	var fromFile string
	cmd := &cobra.Command{
		Use: "ptr-create <vm-id> <ip-id> --from-file <path>", Short: "Create a PTR record",
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := vmID(args[0])
			if err != nil {
				return err
			}
			ip, err := cliutil.ParseInt(args[1])
			if err != nil {
				return err
			}
			body, err := cliutil.ReadBody[api.VPSV1VirtualMachinePTRStoreRequest](fromFile)
			if err != nil {
				return err
			}
			c, err := clictx.FromCommand(cmd).Client()
			if err != nil {
				return err
			}
			resp, err := c.VPSCreatePTRRecordV1WithResponse(cmd.Context(), id, api.IpAddressId(ip), body)
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

func newVMPTRDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use: "ptr-delete <vm-id> <ip-id>", Short: "Delete a PTR record", Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := vmID(args[0])
			if err != nil {
				return err
			}
			ip, err := cliutil.ParseInt(args[1])
			if err != nil {
				return err
			}
			c, err := clictx.FromCommand(cmd).Client()
			if err != nil {
				return err
			}
			resp, err := c.VPSDeletePTRRecordV1WithResponse(cmd.Context(), id, api.IpAddressId(ip))
			if err != nil {
				return err
			}
			return clictx.Render(cmd, resp.JSON200)
		},
	}
}

func newVMMonarxInstallCmd() *cobra.Command {
	return &cobra.Command{
		Use: "monarx-install <vm-id>", Short: "Install Monarx malware scanner", Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := vmID(args[0])
			if err != nil {
				return err
			}
			c, err := clictx.FromCommand(cmd).Client()
			if err != nil {
				return err
			}
			resp, err := c.VPSInstallMonarxV1WithResponse(cmd.Context(), id)
			if err != nil {
				return err
			}
			return clictx.Render(cmd, resp.JSON200)
		},
	}
}

func newVMMonarxUninstallCmd() *cobra.Command {
	return &cobra.Command{
		Use: "monarx-uninstall <vm-id>", Short: "Uninstall Monarx", Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := vmID(args[0])
			if err != nil {
				return err
			}
			c, err := clictx.FromCommand(cmd).Client()
			if err != nil {
				return err
			}
			resp, err := c.VPSUninstallMonarxV1WithResponse(cmd.Context(), id)
			if err != nil {
				return err
			}
			return clictx.Render(cmd, resp.JSON200)
		},
	}
}

func newVMMonarxMetricsCmd() *cobra.Command {
	return &cobra.Command{
		Use: "monarx-metrics <vm-id>", Short: "Get Monarx scan metrics", Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := vmID(args[0])
			if err != nil {
				return err
			}
			c, err := clictx.FromCommand(cmd).Client()
			if err != nil {
				return err
			}
			resp, err := c.VPSGetScanMetricsV1WithResponse(cmd.Context(), id)
			if err != nil {
				return err
			}
			return clictx.Render(cmd, resp.JSON200)
		},
	}
}
