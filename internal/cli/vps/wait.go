package vps

import (
	"time"

	"github.com/spf13/cobra"

	"github.com/rizaleow/hostinger-cli/internal/api"
	"github.com/rizaleow/hostinger-cli/internal/clictx"
	"github.com/rizaleow/hostinger-cli/internal/poll"
)

// addWaitFlags adds --wait/--wait-timeout flags to a command whose RunE
// returns an action resource.
func addWaitFlags(cmd *cobra.Command) {
	cmd.Flags().Bool("wait", false, "block until the action reaches a terminal state")
	cmd.Flags().Duration("wait-timeout", 5*time.Minute, "max time to wait when --wait is set")
}

// maybeWait inspects the action and, if --wait was passed, polls until it
// reaches a terminal state before rendering.
func maybeWait(cmd *cobra.Command, vm api.VirtualMachineId, action *api.VPSV1ActionActionResource) error {
	wait, _ := cmd.Flags().GetBool("wait")
	timeout, _ := cmd.Flags().GetDuration("wait-timeout")
	if !wait || action == nil || action.Id == nil {
		return clictx.Render(cmd, action)
	}
	c, err := clictx.FromCommand(cmd).Client()
	if err != nil {
		return err
	}
	final, err := poll.Wait(cmd.Context(), vm, *action.Id, poll.FetchVMAction(c),
		poll.Options{Timeout: timeout})
	if err != nil {
		return err
	}
	return clictx.Render(cmd, final)
}
