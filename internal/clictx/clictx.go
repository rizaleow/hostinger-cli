// Package clictx is the import-cycle-free bridge between the root command
// (which builds the runtime State) and individual subcommand packages
// (which read it). The root attaches a *State to the cobra command's
// context; subcommands fetch it via FromCommand.
package clictx

import (
	"context"
	"fmt"
	"io"

	"github.com/spf13/cobra"

	"github.com/rizaleow/hostinger-cli/internal/api"
	"github.com/rizaleow/hostinger-cli/internal/config"
	"github.com/rizaleow/hostinger-cli/internal/output"
)

type ctxKey int

const stateKey ctxKey = 0

// State is the runtime context handed to every subcommand.
type State struct {
	Config        *config.File
	ConfigPath    string
	Resolved      config.Resolved
	OutputOptions output.Options

	clientErr error
	client    *api.ClientWithResponses
}

// Client lazily constructs the typed Hostinger client. Returns an error if
// no token was resolved.
func (s *State) Client() (*api.ClientWithResponses, error) {
	if s.client != nil || s.clientErr != nil {
		return s.client, s.clientErr
	}
	if s.Resolved.Token == "" {
		s.clientErr = fmt.Errorf("no API token: set %s, run `hostinger-cli auth login`, or pass --token", config.EnvToken)
		return nil, s.clientErr
	}
	s.client, s.clientErr = api.New(api.Options{
		BaseURL: s.Resolved.BaseURL,
		Token:   s.Resolved.Token,
	})
	return s.client, s.clientErr
}

// With attaches the state to ctx.
func With(ctx context.Context, s *State) context.Context {
	return context.WithValue(ctx, stateKey, s)
}

// FromContext extracts the state, or returns nil if absent.
func FromContext(ctx context.Context) *State {
	s, _ := ctx.Value(stateKey).(*State)
	return s
}

// FromCommand is the canonical accessor used inside RunE bodies.
func FromCommand(cmd *cobra.Command) *State {
	s := FromContext(cmd.Context())
	if s == nil {
		panic("hostinger-cli: clictx.State missing — root PersistentPreRunE did not run")
	}
	return s
}

// Render writes payload using the resolved output options on the cmd's state.
func Render(cmd *cobra.Command, payload any) error {
	s := FromCommand(cmd)
	opts := s.OutputOptions
	if opts.Out == nil {
		opts.Out = cmd.OutOrStdout()
	}
	return output.Render(payload, opts)
}

// Out returns the configured output writer.
func Out(cmd *cobra.Command) io.Writer { return cmd.OutOrStdout() }
