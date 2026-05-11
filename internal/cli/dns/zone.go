package dns

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/rizaleow/hostinger-cli/internal/api"
	"github.com/rizaleow/hostinger-cli/internal/clictx"
)

func newZoneCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "zone", Short: "DNS zone records"}
	cmd.AddCommand(
		newZoneGetCmd(),
		newZoneUpdateCmd(),
		newZoneDeleteCmd(),
		newZoneResetCmd(),
		newZoneValidateCmd(),
	)
	return cmd
}

func newZoneGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <domain>",
		Short: "Get DNS records for a domain",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clictx.FromCommand(cmd).Client()
			if err != nil {
				return err
			}
			resp, err := client.DNSGetDNSRecordsV1WithResponse(cmd.Context(), api.Domain(args[0]))
			if err != nil {
				return err
			}
			return clictx.Render(cmd, resp.JSON200)
		},
	}
}

func newZoneUpdateCmd() *cobra.Command {
	var bodyFile string
	cmd := &cobra.Command{
		Use:   "update <domain> --from-file <path>",
		Short: "Update DNS records from a JSON request body",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			body, err := readBody[api.DNSV1ZoneUpdateRequest](bodyFile)
			if err != nil {
				return err
			}
			client, err := clictx.FromCommand(cmd).Client()
			if err != nil {
				return err
			}
			resp, err := client.DNSUpdateDNSRecordsV1WithResponse(cmd.Context(), api.Domain(args[0]), body)
			if err != nil {
				return err
			}
			return clictx.Render(cmd, resp.JSON200)
		},
	}
	cmd.Flags().StringVar(&bodyFile, "from-file", "", "path to JSON file with update body ('-' for stdin)")
	_ = cmd.MarkFlagRequired("from-file")
	return cmd
}

func newZoneDeleteCmd() *cobra.Command {
	var bodyFile string
	cmd := &cobra.Command{
		Use:   "delete <domain> --from-file <path>",
		Short: "Delete DNS records described in a JSON request body",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			body, err := readBody[api.DNSV1ZoneDestroyRequest](bodyFile)
			if err != nil {
				return err
			}
			client, err := clictx.FromCommand(cmd).Client()
			if err != nil {
				return err
			}
			resp, err := client.DNSDeleteDNSRecordsV1WithResponse(cmd.Context(), api.Domain(args[0]), body)
			if err != nil {
				return err
			}
			return clictx.Render(cmd, resp.JSON200)
		},
	}
	cmd.Flags().StringVar(&bodyFile, "from-file", "", "path to JSON file with delete filter body ('-' for stdin)")
	_ = cmd.MarkFlagRequired("from-file")
	return cmd
}

func newZoneResetCmd() *cobra.Command {
	var bodyFile string
	cmd := &cobra.Command{
		Use:   "reset <domain> --from-file <path>",
		Short: "Reset DNS zone to defaults",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			body, err := readBody[api.DNSV1ZoneResetRequest](bodyFile)
			if err != nil {
				return err
			}
			client, err := clictx.FromCommand(cmd).Client()
			if err != nil {
				return err
			}
			resp, err := client.DNSResetDNSRecordsV1WithResponse(cmd.Context(), api.Domain(args[0]), body)
			if err != nil {
				return err
			}
			return clictx.Render(cmd, resp.JSON200)
		},
	}
	cmd.Flags().StringVar(&bodyFile, "from-file", "", "path to JSON file with reset body ('-' for stdin)")
	_ = cmd.MarkFlagRequired("from-file")
	return cmd
}

func newZoneValidateCmd() *cobra.Command {
	var bodyFile string
	cmd := &cobra.Command{
		Use:   "validate <domain> --from-file <path>",
		Short: "Validate DNS records without applying changes",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			body, err := readBody[api.DNSV1ZoneUpdateRequest](bodyFile)
			if err != nil {
				return err
			}
			client, err := clictx.FromCommand(cmd).Client()
			if err != nil {
				return err
			}
			resp, err := client.DNSValidateDNSRecordsV1WithResponse(cmd.Context(), api.Domain(args[0]), body)
			if err != nil {
				return err
			}
			return clictx.Render(cmd, resp.JSON200)
		},
	}
	cmd.Flags().StringVar(&bodyFile, "from-file", "", "path to JSON file with body ('-' for stdin)")
	_ = cmd.MarkFlagRequired("from-file")
	return cmd
}

// readBody decodes JSON from a file path or stdin into T.
func readBody[T any](path string) (T, error) {
	var v T
	var data []byte
	var err error
	if path == "-" {
		data, err = readAll(os.Stdin)
	} else {
		data, err = os.ReadFile(path)
	}
	if err != nil {
		return v, fmt.Errorf("read body: %w", err)
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return v, fmt.Errorf("parse body JSON: %w", err)
	}
	return v, nil
}

func readAll(f *os.File) ([]byte, error) {
	buf := make([]byte, 0, 4096)
	tmp := make([]byte, 4096)
	for {
		n, err := f.Read(tmp)
		if n > 0 {
			buf = append(buf, tmp[:n]...)
		}
		if err != nil {
			if err.Error() == "EOF" {
				return buf, nil
			}
			return buf, err
		}
	}
}
