package magalucloud

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/MagaluCloud/mgc-sdk-go/compute"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/hashicorp/packer-plugin-sdk/packer"
)

type StepWaitInstanceTeardown struct{}

func (s *StepWaitInstanceTeardown) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	ui := state.Get("ui").(packer.Ui)
	cli := state.Get("compute").(*compute.VirtualMachineClient)
	id := state.Get("instance_id").(string)

	ui.Say(fmt.Sprintf("Waiting for virtual machine instance teardown with ID: %s", id))

	ticker := time.NewTicker(waitInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			state.Put("error", ctx.Err())
			return multistep.ActionHalt
		case <-ticker.C:
			instance, err := cli.Instances().Get(ctx, id, []compute.InstanceExpand{})
			if err != nil && strings.Contains(err.Error(), "404") {
				return multistep.ActionContinue
			}
			if err != nil {
				state.Put("error", fmt.Errorf("Error querying virtual machine: %s", err))
				return multistep.ActionHalt
			}
			if instance.State == "error" {
				state.Put("error", fmt.Errorf("Virtual machine state error: %s", instance.Error.Message))
				return multistep.ActionHalt
			}
			if strings.Contains(instance.Status, "error") {
				state.Put("error", fmt.Errorf("Virtual machine status error: %s", instance.Error.Message))
				return multistep.ActionHalt
			}
		}
	}
}

func (s *StepWaitInstanceTeardown) Cleanup(_ multistep.StateBag) {}
