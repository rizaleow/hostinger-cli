package vps

import (
	"github.com/spf13/cobra"

	"github.com/rizaleow/hostinger-cli/internal/api"
	"github.com/rizaleow/hostinger-cli/internal/clictx"
	"github.com/rizaleow/hostinger-cli/internal/cliutil"
)

func newFirewallCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "firewall", Short: "VPS firewalls and rules"}
	cmd.AddCommand(
		newFWListCmd(), newFWGetCmd(), newFWCreateCmd(), newFWDeleteCmd(),
		newFWActivateCmd(), newFWDeactivateCmd(), newFWSyncCmd(),
		newFWRuleCreateCmd(), newFWRuleUpdateCmd(), newFWRuleDeleteCmd(),
	)
	return cmd
}

func newFWListCmd() *cobra.Command {
	return &cobra.Command{
		Use: "list", Short: "List firewalls",
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := clictx.FromCommand(cmd).Client()
			if err != nil {
				return err
			}
			resp, err := c.VPSGetFirewallListV1WithResponse(cmd.Context(), &api.VPSGetFirewallListV1Params{})
			if err != nil {
				return err
			}
			return clictx.Render(cmd, resp.JSON200)
		},
	}
}

func fwID(arg string) (api.FirewallId, error) {
	n, err := cliutil.ParseInt(arg)
	if err != nil {
		return 0, err
	}
	return api.FirewallId(n), nil
}

func newFWGetCmd() *cobra.Command {
	return &cobra.Command{
		Use: "get <firewall-id>", Short: "Get firewall details", Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := fwID(args[0])
			if err != nil {
				return err
			}
			c, err := clictx.FromCommand(cmd).Client()
			if err != nil {
				return err
			}
			resp, err := c.VPSGetFirewallDetailsV1WithResponse(cmd.Context(), id)
			if err != nil {
				return err
			}
			return clictx.Render(cmd, resp.JSON200)
		},
	}
}

func newFWCreateCmd() *cobra.Command {
	var fromFile string
	cmd := &cobra.Command{
		Use: "create --from-file <path>", Short: "Create a new firewall",
		RunE: func(cmd *cobra.Command, _ []string) error {
			body, err := cliutil.ReadBody[api.VPSV1FirewallStoreRequest](fromFile)
			if err != nil {
				return err
			}
			c, err := clictx.FromCommand(cmd).Client()
			if err != nil {
				return err
			}
			resp, err := c.VPSCreateNewFirewallV1WithResponse(cmd.Context(), body)
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

func newFWDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use: "delete <firewall-id>", Short: "Delete a firewall", Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := fwID(args[0])
			if err != nil {
				return err
			}
			c, err := clictx.FromCommand(cmd).Client()
			if err != nil {
				return err
			}
			resp, err := c.VPSDeleteFirewallV1WithResponse(cmd.Context(), id)
			if err != nil {
				return err
			}
			return clictx.Render(cmd, resp.JSON200)
		},
	}
}

func newFWActivateCmd() *cobra.Command {
	return &cobra.Command{
		Use: "activate <firewall-id> <vm-id>", Short: "Attach firewall to a VM", Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			fid, err := fwID(args[0])
			if err != nil {
				return err
			}
			vid, err := vmID(args[1])
			if err != nil {
				return err
			}
			c, err := clictx.FromCommand(cmd).Client()
			if err != nil {
				return err
			}
			resp, err := c.VPSActivateFirewallV1WithResponse(cmd.Context(), fid, vid)
			if err != nil {
				return err
			}
			return clictx.Render(cmd, resp.JSON200)
		},
	}
}

func newFWDeactivateCmd() *cobra.Command {
	return &cobra.Command{
		Use: "deactivate <firewall-id> <vm-id>", Short: "Detach firewall from a VM", Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			fid, err := fwID(args[0])
			if err != nil {
				return err
			}
			vid, err := vmID(args[1])
			if err != nil {
				return err
			}
			c, err := clictx.FromCommand(cmd).Client()
			if err != nil {
				return err
			}
			resp, err := c.VPSDeactivateFirewallV1WithResponse(cmd.Context(), fid, vid)
			if err != nil {
				return err
			}
			return clictx.Render(cmd, resp.JSON200)
		},
	}
}

func newFWSyncCmd() *cobra.Command {
	return &cobra.Command{
		Use: "sync <firewall-id> <vm-id>", Short: "Sync firewall rules to a VM", Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			fid, err := fwID(args[0])
			if err != nil {
				return err
			}
			vid, err := vmID(args[1])
			if err != nil {
				return err
			}
			c, err := clictx.FromCommand(cmd).Client()
			if err != nil {
				return err
			}
			resp, err := c.VPSSyncFirewallV1WithResponse(cmd.Context(), fid, vid)
			if err != nil {
				return err
			}
			return clictx.Render(cmd, resp.JSON200)
		},
	}
}

func newFWRuleCreateCmd() *cobra.Command {
	var fromFile string
	cmd := &cobra.Command{
		Use: "rule-create <firewall-id> --from-file <path>", Short: "Add a firewall rule", Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			fid, err := fwID(args[0])
			if err != nil {
				return err
			}
			body, err := cliutil.ReadBody[api.VPSV1FirewallRulesStoreRequest](fromFile)
			if err != nil {
				return err
			}
			c, err := clictx.FromCommand(cmd).Client()
			if err != nil {
				return err
			}
			resp, err := c.VPSCreateFirewallRuleV1WithResponse(cmd.Context(), fid, body)
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

func newFWRuleUpdateCmd() *cobra.Command {
	var fromFile string
	cmd := &cobra.Command{
		Use: "rule-update <firewall-id> <rule-id> --from-file <path>", Short: "Update a firewall rule",
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			fid, err := fwID(args[0])
			if err != nil {
				return err
			}
			rid, err := cliutil.ParseInt(args[1])
			if err != nil {
				return err
			}
			body, err := cliutil.ReadBody[api.VPSV1FirewallRulesStoreRequest](fromFile)
			if err != nil {
				return err
			}
			c, err := clictx.FromCommand(cmd).Client()
			if err != nil {
				return err
			}
			resp, err := c.VPSUpdateFirewallRuleV1WithResponse(cmd.Context(), fid, api.RuleId(rid), body)
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

func newFWRuleDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use: "rule-delete <firewall-id> <rule-id>", Short: "Delete a firewall rule", Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			fid, err := fwID(args[0])
			if err != nil {
				return err
			}
			rid, err := cliutil.ParseInt(args[1])
			if err != nil {
				return err
			}
			c, err := clictx.FromCommand(cmd).Client()
			if err != nil {
				return err
			}
			resp, err := c.VPSDeleteFirewallRuleV1WithResponse(cmd.Context(), fid, api.RuleId(rid))
			if err != nil {
				return err
			}
			return clictx.Render(cmd, resp.JSON200)
		},
	}
}
