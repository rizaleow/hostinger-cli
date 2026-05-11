// Package config resolves authentication tokens and runtime settings from
// (in precedence order): explicit flags, environment variables, OS keychain,
// and on-disk YAML config.
package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// EnvToken is the environment variable used to override the bearer token.
const EnvToken = "HOSTINGER_API_TOKEN"

// EnvOutput overrides the default output format when no -o flag is given.
const EnvOutput = "HOSTINGER_OUTPUT"

// EnvBaseURL overrides the Hostinger API base URL (useful for testing).
const EnvBaseURL = "HOSTINGER_API_BASE_URL"

// Profile holds the persisted settings for a single named profile.
type Profile struct {
	Token   string `yaml:"token,omitempty"`
	BaseURL string `yaml:"base_url,omitempty"`
}

// File is the on-disk YAML config schema.
type File struct {
	CurrentProfile string             `yaml:"current_profile,omitempty"`
	UseKeyring     bool               `yaml:"use_keyring,omitempty"`
	Profiles       map[string]Profile `yaml:"profiles,omitempty"`
}

// DefaultPath returns ~/.config/hostinger-cli/config.yaml (XDG-aware).
func DefaultPath() (string, error) {
	if env := os.Getenv("HOSTINGER_CLI_CONFIG"); env != "" {
		return env, nil
	}
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		return filepath.Join(xdg, "hostinger-cli", "config.yaml"), nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "hostinger-cli", "config.yaml"), nil
}

// Load reads the config file from path. A missing file is not an error.
func Load(path string) (*File, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &File{Profiles: map[string]Profile{}}, nil
		}
		return nil, err
	}
	var f File
	if err := yaml.Unmarshal(data, &f); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}
	if f.Profiles == nil {
		f.Profiles = map[string]Profile{}
	}
	return &f, nil
}

// Save writes the config file (creating parent dirs and chmod 0600).
func Save(path string, f *File) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	data, err := yaml.Marshal(f)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o600)
}

// Profile returns the profile by name, or an empty profile if absent.
func (f *File) Profile(name string) Profile {
	if name == "" {
		name = f.CurrentProfile
	}
	if name == "" {
		name = "default"
	}
	return f.Profiles[name]
}

// SetProfile upserts a profile.
func (f *File) SetProfile(name string, p Profile) {
	if name == "" {
		name = "default"
	}
	if f.Profiles == nil {
		f.Profiles = map[string]Profile{}
	}
	f.Profiles[name] = p
	if f.CurrentProfile == "" {
		f.CurrentProfile = name
	}
}
