// Package reach implements `hostinger-cli reach` (email marketing).
package reach

import (
	"github.com/google/uuid"
	"github.com/spf13/cobra"

	"github.com/rizaleow/hostinger-cli/internal/api"
	"github.com/rizaleow/hostinger-cli/internal/clictx"
	"github.com/rizaleow/hostinger-cli/internal/cliutil"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "reach", Short: "Reach: email marketing contacts, segments, profiles"}
	cmd.AddCommand(newContactsCmd(), newSegmentsCmd(), newProfilesCmd())
	return cmd
}

func newContactsCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "contacts", Short: "Email contacts"}
	cmd.AddCommand(
		&cobra.Command{
			Use: "list", Short: "List contacts",
			RunE: func(cmd *cobra.Command, _ []string) error {
				c, err := clictx.FromCommand(cmd).Client()
				if err != nil {
					return err
				}
				params := &api.ReachListContactsV1Params{}
				resp, err := c.ReachListContactsV1WithResponse(cmd.Context(), params)
				if err != nil {
					return err
				}
				return clictx.Render(cmd, resp.JSON200)
			},
		},
		&cobra.Command{
			Use: "delete <uuid>", Short: "Delete a contact", Args: cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				id, err := uuid.Parse(args[0])
				if err != nil {
					return err
				}
				c, err := clictx.FromCommand(cmd).Client()
				if err != nil {
					return err
				}
				resp, err := c.ReachDeleteAContactV1WithResponse(cmd.Context(), api.Uuid(id))
				if err != nil {
					return err
				}
				return clictx.Render(cmd, resp.JSON200)
			},
		},
		newContactsCreateInProfileCmd(),
	)
	return cmd
}

func newContactsCreateInProfileCmd() *cobra.Command {
	var fromFile string
	cmd := &cobra.Command{
		Use: "create <profile-uuid> --from-file <path>", Short: "Create a contact in a profile", Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			body, err := cliutil.ReadBody[api.ReachV1ContactsStoreRequest](fromFile)
			if err != nil {
				return err
			}
			c, err := clictx.FromCommand(cmd).Client()
			if err != nil {
				return err
			}
			resp, err := c.ReachCreateNewContactsV1WithResponse(cmd.Context(), api.ProfileUuid(args[0]), body)
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

func newSegmentsCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "segments", Short: "Contact segments"}
	cmd.AddCommand(
		&cobra.Command{
			Use: "list", Short: "List segments",
			RunE: func(cmd *cobra.Command, _ []string) error {
				c, err := clictx.FromCommand(cmd).Client()
				if err != nil {
					return err
				}
				resp, err := c.ReachListSegmentsV1WithResponse(cmd.Context())
				if err != nil {
					return err
				}
				return clictx.Render(cmd, resp.JSON200)
			},
		},
		&cobra.Command{
			Use: "get <segment-uuid>", Short: "Get a segment", Args: cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				c, err := clictx.FromCommand(cmd).Client()
				if err != nil {
					return err
				}
				resp, err := c.ReachGetSegmentDetailsV1WithResponse(cmd.Context(), api.SegmentUuid(args[0]))
				if err != nil {
					return err
				}
				return clictx.Render(cmd, resp.JSON200)
			},
		},
		newSegmentsCreateCmd(),
		newSegmentsContactsCmd(),
	)
	return cmd
}

func newSegmentsCreateCmd() *cobra.Command {
	var fromFile string
	cmd := &cobra.Command{
		Use: "create --from-file <path>", Short: "Create a contact segment",
		RunE: func(cmd *cobra.Command, _ []string) error {
			body, err := cliutil.ReadBody[api.ReachV1ContactsSegmentsStoreRequest](fromFile)
			if err != nil {
				return err
			}
			c, err := clictx.FromCommand(cmd).Client()
			if err != nil {
				return err
			}
			resp, err := c.ReachCreateANewContactSegmentV1WithResponse(cmd.Context(), body)
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

func newSegmentsContactsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "list-contacts <segment-uuid>", Short: "List contacts in a segment", Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := clictx.FromCommand(cmd).Client()
			if err != nil {
				return err
			}
			params := &api.ReachListSegmentContactsV1Params{}
			resp, err := c.ReachListSegmentContactsV1WithResponse(cmd.Context(), api.SegmentUuid(args[0]), params)
			if err != nil {
				return err
			}
			return clictx.Render(cmd, resp.JSON200)
		},
	}
	return cmd
}

func newProfilesCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "profiles", Short: "Reach profiles"}
	cmd.AddCommand(&cobra.Command{
		Use: "list", Short: "List profiles",
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := clictx.FromCommand(cmd).Client()
			if err != nil {
				return err
			}
			resp, err := c.ReachListProfilesV1WithResponse(cmd.Context())
			if err != nil {
				return err
			}
			return clictx.Render(cmd, resp.JSON200)
		},
	})
	return cmd
}
