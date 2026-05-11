package vps

import (
	"github.com/spf13/cobra"

	"github.com/rizaleow/hostinger-cli/internal/api"
	"github.com/rizaleow/hostinger-cli/internal/clictx"
	"github.com/rizaleow/hostinger-cli/internal/cliutil"
)

func newPublicKeysCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "public-keys", Short: "SSH public keys"}
	cmd.AddCommand(
		&cobra.Command{
			Use: "list", Short: "List public keys",
			RunE: func(cmd *cobra.Command, _ []string) error {
				c, err := clictx.FromCommand(cmd).Client()
				if err != nil {
					return err
				}
				resp, err := c.VPSGetPublicKeysV1WithResponse(cmd.Context(), &api.VPSGetPublicKeysV1Params{})
				if err != nil {
					return err
				}
				return clictx.Render(cmd, resp.JSON200)
			},
		},
		newPublicKeysCreateCmd(),
		&cobra.Command{
			Use: "delete <id>", Short: "Delete a public key", Args: cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				id, err := cliutil.ParseInt(args[0])
				if err != nil {
					return err
				}
				c, err := clictx.FromCommand(cmd).Client()
				if err != nil {
					return err
				}
				resp, err := c.VPSDeletePublicKeyV1WithResponse(cmd.Context(), api.PublicKeyId(id))
				if err != nil {
					return err
				}
				return clictx.Render(cmd, resp.JSON200)
			},
		},
		newPublicKeysAttachCmd(),
	)
	return cmd
}

func newPublicKeysCreateCmd() *cobra.Command {
	var fromFile string
	cmd := &cobra.Command{
		Use: "create --from-file <path>", Short: "Add a public key",
		RunE: func(cmd *cobra.Command, _ []string) error {
			body, err := cliutil.ReadBody[api.VPSV1PublicKeyStoreRequest](fromFile)
			if err != nil {
				return err
			}
			c, err := clictx.FromCommand(cmd).Client()
			if err != nil {
				return err
			}
			resp, err := c.VPSCreatePublicKeyV1WithResponse(cmd.Context(), body)
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

func newPublicKeysAttachCmd() *cobra.Command {
	var fromFile string
	cmd := &cobra.Command{
		Use: "attach <vm-id> --from-file <path>", Short: "Attach public keys to a VM", Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := vmID(args[0])
			if err != nil {
				return err
			}
			body, err := cliutil.ReadBody[api.VPSV1PublicKeyAttachRequest](fromFile)
			if err != nil {
				return err
			}
			c, err := clictx.FromCommand(cmd).Client()
			if err != nil {
				return err
			}
			resp, err := c.VPSAttachPublicKeyV1WithResponse(cmd.Context(), id, body)
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
