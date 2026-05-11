// Package auth implements the `hostinger-cli auth` command tree.
package auth

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/term"

	"github.com/rizaleow/hostinger-cli/internal/api"
	"github.com/rizaleow/hostinger-cli/internal/config"
)

// stateAccessor avoids an import cycle with the cli package: we only need a
// few fields off cli.State to read the on-disk config.
type stateAccessor interface {
	configFile() *config.File
	configPath() string
}

// NewCmd returns the `auth` command group.
func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth",
		Short: "Manage API authentication",
	}
	cmd.AddCommand(newLoginCmd(), newLogoutCmd(), newStatusCmd(), newWhoamiCmd())
	return cmd
}

func newLoginCmd() *cobra.Command {
	var (
		tokenFlag  string
		useKeyring bool
		profile    string
	)
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Store an API token for future commands",
		Long: "Prompts for an API token (or reads --token / stdin), validates it against the Hostinger API, " +
			"and saves it to the config file. With --keyring, the token is stored in your OS keychain instead " +
			"of in plaintext.",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			path, err := loadPath(cmd)
			if err != nil {
				return err
			}
			file, err := config.Load(path)
			if err != nil {
				return err
			}
			if profile == "" {
				profile = file.CurrentProfile
				if profile == "" {
					profile = "default"
				}
			}

			tok := strings.TrimSpace(tokenFlag)
			if tok == "" {
				tok, err = promptToken(cmd)
				if err != nil {
					return err
				}
			}
			if tok == "" {
				return errors.New("no token provided")
			}

			if err := validateToken(ctx, tok); err != nil {
				return fmt.Errorf("token validation failed: %w", err)
			}

			p := file.Profile(profile)
			if useKeyring {
				if err := config.KeyringSet(profile, tok); err != nil {
					return fmt.Errorf("keyring: %w", err)
				}
				p.Token = ""
				file.UseKeyring = true
			} else {
				p.Token = tok
			}
			file.SetProfile(profile, p)
			file.CurrentProfile = profile
			if err := config.Save(path, file); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Logged in to profile %q. Token stored in %s.\n",
				profile, storageLocation(useKeyring, path))
			return nil
		},
	}
	cmd.Flags().StringVar(&tokenFlag, "token", "", "API token (skips prompt)")
	cmd.Flags().BoolVar(&useKeyring, "keyring", false, "store token in OS keychain instead of config file")
	cmd.Flags().StringVar(&profile, "as-profile", "", "name of the profile to write (default: current or 'default')")
	cmd.Annotations = map[string]string{"skip-state": "true"}
	return cmd
}

func newLogoutCmd() *cobra.Command {
	var profile string
	cmd := &cobra.Command{
		Use:   "logout",
		Short: "Remove a stored API token",
		RunE: func(cmd *cobra.Command, args []string) error {
			path, err := loadPath(cmd)
			if err != nil {
				return err
			}
			file, err := config.Load(path)
			if err != nil {
				return err
			}
			if profile == "" {
				profile = file.CurrentProfile
				if profile == "" {
					profile = "default"
				}
			}
			p := file.Profile(profile)
			p.Token = ""
			file.SetProfile(profile, p)
			_ = config.KeyringDelete(profile)
			if err := config.Save(path, file); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Logged out of profile %q.\n", profile)
			return nil
		},
	}
	cmd.Flags().StringVar(&profile, "as-profile", "", "profile to clear")
	cmd.Annotations = map[string]string{"skip-state": "true"}
	return cmd
}

func newStatusCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Show which token (if any) is currently active",
		RunE: func(cmd *cobra.Command, args []string) error {
			path, err := loadPath(cmd)
			if err != nil {
				return err
			}
			file, err := config.Load(path)
			if err != nil {
				return err
			}
			r := config.Resolve(file, config.ResolveOptions{
				Profile: cmd.Flag("profile").Value.String(),
			})
			out := cmd.OutOrStdout()
			if r.Token == "" {
				fmt.Fprintln(out, "Not logged in.")
				fmt.Fprintf(out, "Try: hostinger-cli auth login\n")
				return nil
			}
			fmt.Fprintf(out, "Profile: %s\n", r.Profile)
			fmt.Fprintf(out, "Source:  %s\n", r.Source)
			fmt.Fprintf(out, "Token:   %s\n", redact(r.Token))
			if r.BaseURL != "" {
				fmt.Fprintf(out, "BaseURL: %s\n", r.BaseURL)
			}
			return nil
		},
	}
	cmd.Annotations = map[string]string{"skip-state": "true"}
	return cmd
}

func newWhoamiCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "whoami",
		Short: "Validate the active token against the API",
		RunE: func(cmd *cobra.Command, args []string) error {
			path, err := loadPath(cmd)
			if err != nil {
				return err
			}
			file, err := config.Load(path)
			if err != nil {
				return err
			}
			r := config.Resolve(file, config.ResolveOptions{
				Profile: cmd.Flag("profile").Value.String(),
			})
			if r.Token == "" {
				return errors.New("no API token: run `hostinger-cli auth login`")
			}
			if err := validateToken(cmd.Context(), r.Token); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Token for profile %q is valid.\n", r.Profile)
			return nil
		},
	}
	cmd.Annotations = map[string]string{"skip-state": "true"}
	return cmd
}

// validateToken performs a cheap authenticated request to confirm the token works.
func validateToken(ctx context.Context, tok string) error {
	c, err := api.New(api.Options{Token: tok})
	if err != nil {
		return err
	}
	resp, err := c.BillingGetPaymentMethodListV1WithResponse(ctx)
	if err != nil {
		return err
	}
	switch resp.StatusCode() {
	case http.StatusOK:
		return nil
	case http.StatusUnauthorized:
		return errors.New("unauthorized: token rejected by API")
	default:
		return fmt.Errorf("unexpected status %d", resp.StatusCode())
	}
}

func promptToken(cmd *cobra.Command) (string, error) {
	in := cmd.InOrStdin()
	if f, ok := in.(*os.File); ok && term.IsTerminal(int(f.Fd())) {
		fmt.Fprint(cmd.ErrOrStderr(), "Hostinger API token: ")
		b, err := term.ReadPassword(int(f.Fd()))
		fmt.Fprintln(cmd.ErrOrStderr())
		if err != nil {
			return "", err
		}
		return strings.TrimSpace(string(b)), nil
	}
	// Non-interactive: read a line from stdin.
	r := bufio.NewReader(in)
	line, err := r.ReadString('\n')
	if err != nil && !errors.Is(err, errReadEOF()) {
		return "", err
	}
	return strings.TrimSpace(line), nil
}

func errReadEOF() error { return errors.New("EOF") }

func redact(t string) string {
	if len(t) <= 8 {
		return strings.Repeat("•", len(t))
	}
	return t[:4] + strings.Repeat("•", 12) + t[len(t)-4:]
}

func storageLocation(keyring bool, path string) string {
	if keyring {
		return "OS keychain"
	}
	return path
}

func loadPath(cmd *cobra.Command) (string, error) {
	if f := cmd.Flag("config"); f != nil && f.Value.String() != "" {
		return f.Value.String(), nil
	}
	return config.DefaultPath()
}
