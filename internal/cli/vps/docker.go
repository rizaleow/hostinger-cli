package vps

import (
	"context"
	"time"

	"github.com/spf13/cobra"

	"github.com/rizaleow/hostinger-cli/internal/api"
	"github.com/rizaleow/hostinger-cli/internal/clictx"
	"github.com/rizaleow/hostinger-cli/internal/cliutil"
	"github.com/rizaleow/hostinger-cli/internal/poll"
)

type projectOp func(*api.ClientWithResponses, context.Context, api.VirtualMachineId, api.ProjectName) (any, error)

// supportsWait identifies which project ops return an action and can therefore
// be polled with --wait. The other endpoints (containers/logs/get/list) return
// plain data and have no async lifecycle.
var supportsWait = map[string]bool{
	"delete": true, "start": true, "stop": true,
	"restart": true, "update": true,
}

func newDockerCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "docker", Short: "Docker Compose Manager on a VPS (experimental)"}
	cmd.AddCommand(newProjectCmd())
	return cmd
}

func newProjectCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "project", Short: "Docker Compose projects"}
	cmd.AddCommand(
		newProjectListCmd(),
		newProjectGetCmd(),
		newProjectCreateCmd(),
		projectActionCmd("delete", "Delete a Docker project",
			func(c *api.ClientWithResponses, ctx context.Context, v api.VirtualMachineId, p api.ProjectName) (any, error) {
				resp, err := c.VPSDeleteProjectV1WithResponse(ctx, v, p)
				if err != nil {
					return nil, err
				}
				return resp.JSON200, nil
			}),
		projectActionCmd("start", "Start a Docker project",
			func(c *api.ClientWithResponses, ctx context.Context, v api.VirtualMachineId, p api.ProjectName) (any, error) {
				resp, err := c.VPSStartProjectV1WithResponse(ctx, v, p)
				if err != nil {
					return nil, err
				}
				return resp.JSON200, nil
			}),
		projectActionCmd("stop", "Stop a Docker project",
			func(c *api.ClientWithResponses, ctx context.Context, v api.VirtualMachineId, p api.ProjectName) (any, error) {
				resp, err := c.VPSStopProjectV1WithResponse(ctx, v, p)
				if err != nil {
					return nil, err
				}
				return resp.JSON200, nil
			}),
		projectActionCmd("restart", "Restart a Docker project",
			func(c *api.ClientWithResponses, ctx context.Context, v api.VirtualMachineId, p api.ProjectName) (any, error) {
				resp, err := c.VPSRestartProjectV1WithResponse(ctx, v, p)
				if err != nil {
					return nil, err
				}
				return resp.JSON200, nil
			}),
		projectActionCmd("update", "Update a Docker project (pull + recreate)",
			func(c *api.ClientWithResponses, ctx context.Context, v api.VirtualMachineId, p api.ProjectName) (any, error) {
				resp, err := c.VPSUpdateProjectV1WithResponse(ctx, v, p)
				if err != nil {
					return nil, err
				}
				return resp.JSON200, nil
			}),
		projectActionCmd("containers", "List project containers",
			func(c *api.ClientWithResponses, ctx context.Context, v api.VirtualMachineId, p api.ProjectName) (any, error) {
				resp, err := c.VPSGetProjectContainersV1WithResponse(ctx, v, p)
				if err != nil {
					return nil, err
				}
				return resp.JSON200, nil
			}),
		projectActionCmd("logs", "Show recent logs across the project",
			func(c *api.ClientWithResponses, ctx context.Context, v api.VirtualMachineId, p api.ProjectName) (any, error) {
				resp, err := c.VPSGetProjectLogsV1WithResponse(ctx, v, p)
				if err != nil {
					return nil, err
				}
				return resp.JSON200, nil
			}),
	)
	return cmd
}

func newProjectListCmd() *cobra.Command {
	return &cobra.Command{
		Use: "list <vm-id>", Short: "List Docker Compose projects", Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := vmID(args[0])
			if err != nil {
				return err
			}
			c, err := clictx.FromCommand(cmd).Client()
			if err != nil {
				return err
			}
			resp, err := c.VPSGetProjectListV1WithResponse(cmd.Context(), id)
			if err != nil {
				return err
			}
			return clictx.Render(cmd, resp.JSON200)
		},
	}
}

func newProjectGetCmd() *cobra.Command {
	return &cobra.Command{
		Use: "get <vm-id> <project>", Short: "Get a Docker project", Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := vmID(args[0])
			if err != nil {
				return err
			}
			c, err := clictx.FromCommand(cmd).Client()
			if err != nil {
				return err
			}
			resp, err := c.VPSGetProjectContentsV1WithResponse(cmd.Context(), id, api.ProjectName(args[1]))
			if err != nil {
				return err
			}
			return clictx.Render(cmd, resp.JSON200)
		},
	}
}

func newProjectCreateCmd() *cobra.Command {
	var fromFile string
	cmd := &cobra.Command{
		Use: "create <vm-id> --from-file <path>", Short: "Create or replace a Docker project", Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := vmID(args[0])
			if err != nil {
				return err
			}
			body, err := cliutil.ReadBody[api.VPSV1VirtualMachineDockerManagerUpRequest](fromFile)
			if err != nil {
				return err
			}
			c, err := clictx.FromCommand(cmd).Client()
			if err != nil {
				return err
			}
			resp, err := c.VPSCreateNewProjectV1WithResponse(cmd.Context(), id, body)
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

func projectActionCmd(use, short string, op projectOp) *cobra.Command {
	cmd := &cobra.Command{
		Use: use + " <vm-id> <project>", Short: short, Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := vmID(args[0])
			if err != nil {
				return err
			}
			c, err := clictx.FromCommand(cmd).Client()
			if err != nil {
				return err
			}
			out, err := op(c, cmd.Context(), id, api.ProjectName(args[1]))
			if err != nil {
				return err
			}

			wait, _ := cmd.Flags().GetBool("wait")
			timeout, _ := cmd.Flags().GetDuration("wait-timeout")
			if wait {
				if action, ok := out.(*api.VPSV1ActionActionResource); ok && action != nil && action.Id != nil {
					final, err := poll.Wait(cmd.Context(), id, *action.Id,
						poll.FetchVMAction(c), poll.Options{Timeout: timeout})
					if err != nil {
						return err
					}
					return clictx.Render(cmd, final)
				}
			}
			return clictx.Render(cmd, out)
		},
	}
	if supportsWait[use] {
		cmd.Flags().Bool("wait", false, "block until the action reaches a terminal state")
		cmd.Flags().Duration("wait-timeout", 5*time.Minute, "max time to wait when --wait is set")
	}
	return cmd
}
