package vps

import (
	"github.com/spf13/cobra"

	"github.com/rizaleow/hostinger-cli/internal/api"
	"github.com/rizaleow/hostinger-cli/internal/clictx"
	"github.com/rizaleow/hostinger-cli/internal/cliutil"
)

func newPostInstallScriptsCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "post-install", Short: "Post-install scripts"}
	cmd.AddCommand(
		&cobra.Command{
			Use: "list", Short: "List post-install scripts",
			RunE: func(cmd *cobra.Command, _ []string) error {
				c, err := clictx.FromCommand(cmd).Client()
				if err != nil {
					return err
				}
				resp, err := c.VPSGetPostInstallScriptsV1WithResponse(cmd.Context(), &api.VPSGetPostInstallScriptsV1Params{})
				if err != nil {
					return err
				}
				return clictx.Render(cmd, resp.JSON200)
			},
		},
		&cobra.Command{
			Use: "get <id>", Short: "Get a script", Args: cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				id, err := cliutil.ParseInt(args[0])
				if err != nil {
					return err
				}
				c, err := clictx.FromCommand(cmd).Client()
				if err != nil {
					return err
				}
				resp, err := c.VPSGetPostInstallScriptV1WithResponse(cmd.Context(), api.PostInstallScriptId(id))
				if err != nil {
					return err
				}
				return clictx.Render(cmd, resp.JSON200)
			},
		},
		newPostInstallCreateCmd(),
		newPostInstallUpdateCmd(),
		&cobra.Command{
			Use: "delete <id>", Short: "Delete a script", Args: cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				id, err := cliutil.ParseInt(args[0])
				if err != nil {
					return err
				}
				c, err := clictx.FromCommand(cmd).Client()
				if err != nil {
					return err
				}
				resp, err := c.VPSDeletePostInstallScriptV1WithResponse(cmd.Context(), api.PostInstallScriptId(id))
				if err != nil {
					return err
				}
				return clictx.Render(cmd, resp.JSON200)
			},
		},
	)
	return cmd
}

func newPostInstallCreateCmd() *cobra.Command {
	var fromFile string
	cmd := &cobra.Command{
		Use: "create --from-file <path>", Short: "Create a post-install script",
		RunE: func(cmd *cobra.Command, _ []string) error {
			body, err := cliutil.ReadBody[api.VPSV1PostInstallScriptStoreRequest](fromFile)
			if err != nil {
				return err
			}
			c, err := clictx.FromCommand(cmd).Client()
			if err != nil {
				return err
			}
			resp, err := c.VPSCreatePostInstallScriptV1WithResponse(cmd.Context(), body)
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

func newPostInstallUpdateCmd() *cobra.Command {
	var fromFile string
	cmd := &cobra.Command{
		Use: "update <id> --from-file <path>", Short: "Update a post-install script", Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := cliutil.ParseInt(args[0])
			if err != nil {
				return err
			}
			body, err := cliutil.ReadBody[api.VPSV1PostInstallScriptStoreRequest](fromFile)
			if err != nil {
				return err
			}
			c, err := clictx.FromCommand(cmd).Client()
			if err != nil {
				return err
			}
			resp, err := c.VPSUpdatePostInstallScriptV1WithResponse(cmd.Context(), api.PostInstallScriptId(id), body)
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
