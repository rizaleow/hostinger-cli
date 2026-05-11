// Package poll provides a generic --wait helper that polls a VPS action
// resource until it reaches a terminal state.
package poll

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/rizaleow/hostinger-cli/internal/api"
)

// Terminal states emitted by the Hostinger VPS action API.
const (
	StateSuccess = "success"
	StateError   = "error"
)

// Action is the subset of fields we need from VPSV1ActionActionResource. We
// accept the resource by value to avoid coupling to the generated pointer
// shape.
type Action struct {
	ID    int
	Name  string
	State string
}

// FromResource normalises a generated action resource.
func FromResource(r *api.VPSV1ActionActionResource) Action {
	if r == nil {
		return Action{}
	}
	a := Action{}
	if r.Id != nil {
		a.ID = *r.Id
	}
	if r.Name != nil {
		a.Name = *r.Name
	}
	if r.State != nil {
		a.State = string(*r.State)
	}
	return a
}

// Options controls the polling cadence.
type Options struct {
	Interval time.Duration // initial wait between polls (default 1s, capped at 5s with backoff)
	Timeout  time.Duration // hard deadline; 0 means "use ctx only"
}

// Wait polls fetch until it returns an action in a terminal state, or until
// ctx is cancelled or the timeout fires.
func Wait(ctx context.Context, vmID api.VirtualMachineId, actionID int, fetch func(context.Context, api.VirtualMachineId, int) (Action, error), opts Options) (Action, error) {
	if opts.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, opts.Timeout)
		defer cancel()
	}
	interval := opts.Interval
	if interval <= 0 {
		interval = 1 * time.Second
	}

	for {
		a, err := fetch(ctx, vmID, actionID)
		if err != nil {
			return a, err
		}
		switch a.State {
		case StateSuccess:
			return a, nil
		case StateError:
			return a, fmt.Errorf("action %d (%s) failed", a.ID, a.Name)
		}

		select {
		case <-ctx.Done():
			return a, ctx.Err()
		case <-time.After(interval):
		}
		if interval < 5*time.Second {
			interval = min(interval*2, 5*time.Second)
		}
	}
}

// FetchVMAction is the canonical fetcher: it calls the action-details endpoint.
func FetchVMAction(client *api.ClientWithResponses) func(context.Context, api.VirtualMachineId, int) (Action, error) {
	return func(ctx context.Context, vm api.VirtualMachineId, id int) (Action, error) {
		resp, err := client.VPSGetActionDetailsV1WithResponse(ctx, vm, api.ActionId(id))
		if err != nil {
			return Action{}, err
		}
		if resp.JSON200 == nil {
			return Action{}, errors.New("empty action response")
		}
		return FromResource(resp.JSON200), nil
	}
}
