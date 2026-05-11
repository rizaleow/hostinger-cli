// Package config exposes `hostinger-cli config` commands for inspecting and
// mutating the on-disk YAML config file.
package config

import (
	"fmt"

	"github.com/spf13/cobra"

	cfg "github.com/rizaleow/hostinger-cli/internal/config"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage hostinger-cli configuration",
	}
	cmd.AddCommand(newPathCmd(), newListCmd(), newSetCmd(), newUnsetCmd())
	return cmd
}

func newPathCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "path",
		Short: "Print the active config file path",
		RunE: func(cmd *cobra.Command, args []string) error {
			path, err := load(cmd)
			if err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), path)
			return nil
		},
	}
	cmd.Annotations = map[string]string{"skip-state": "true"}
	return cmd
}

func newListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List configured profiles",
		RunE: func(cmd *cobra.Command, args []string) error {
			path, err := load(cmd)
			if err != nil {
				return err
			}
			file, err := cfg.Load(path)
			if err != nil {
				return err
			}
			out := cmd.OutOrStdout()
			if len(file.Profiles) == 0 {
				fmt.Fprintln(out, "No profiles configured.")
				return nil
			}
			for name, p := range file.Profiles {
				marker := " "
				if name == file.CurrentProfile {
					marker = "*"
				}
				fmt.Fprintf(out, "%s %-20s base_url=%s token_set=%t\n",
					marker, name, p.BaseURL, p.Token != "")
			}
			return nil
		},
	}
	cmd.Annotations = map[string]string{"skip-state": "true"}
	return cmd
}

func newSetCmd() *cobra.Command {
	var profile string
	cmd := &cobra.Command{
		Use:   "set <key> <value>",
		Short: "Set a config value (key: base_url, current_profile, use_keyring)",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			path, err := load(cmd)
			if err != nil {
				return err
			}
			file, err := cfg.Load(path)
			if err != nil {
				return err
			}
			key, value := args[0], args[1]
			if profile == "" {
				profile = file.CurrentProfile
				if profile == "" {
					profile = "default"
				}
			}
			switch key {
			case "base_url":
				p := file.Profile(profile)
				p.BaseURL = value
				file.SetProfile(profile, p)
			case "current_profile":
				file.CurrentProfile = value
			case "use_keyring":
				file.UseKeyring = value == "true"
			default:
				return fmt.Errorf("unknown key %q", key)
			}
			return cfg.Save(path, file)
		},
	}
	cmd.Flags().StringVar(&profile, "as-profile", "", "target profile (default: current)")
	cmd.Annotations = map[string]string{"skip-state": "true"}
	return cmd
}

func newUnsetCmd() *cobra.Command {
	var profile string
	cmd := &cobra.Command{
		Use:   "unset <profile>",
		Short: "Delete a named profile",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path, err := load(cmd)
			if err != nil {
				return err
			}
			file, err := cfg.Load(path)
			if err != nil {
				return err
			}
			delete(file.Profiles, args[0])
			if file.CurrentProfile == args[0] {
				file.CurrentProfile = ""
			}
			_ = cfg.KeyringDelete(args[0])
			return cfg.Save(path, file)
		},
	}
	cmd.Flags().StringVar(&profile, "as-profile", "", "alias")
	cmd.Annotations = map[string]string{"skip-state": "true"}
	return cmd
}

func load(cmd *cobra.Command) (string, error) {
	if f := cmd.Flag("config"); f != nil && f.Value.String() != "" {
		return f.Value.String(), nil
	}
	return cfg.DefaultPath()
}
