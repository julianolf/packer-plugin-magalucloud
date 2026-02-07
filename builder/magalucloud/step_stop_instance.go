package magalucloud

import (
	"context"
	"fmt"

	"github.com/MagaluCloud/mgc-sdk-go/compute"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/hashicorp/packer-plugin-sdk/packer"
)

type StepStopInstance struct{}

func (s *StepStopInstance) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	ui := state.Get("ui").(packer.Ui)
	cli := state.Get("compute").(*compute.VirtualMachineClient)
	id := state.Get("instance_id").(string)

	ui.Say(fmt.Sprintf("Stopping virtual machine instance with ID: %s", id))

	err := cli.Instances().Stop(ctx, id)
	if err != nil {
		state.Put("error", fmt.Errorf("Error stopping virtual machine: %s", err))
		return multistep.ActionHalt
	}

	return multistep.ActionContinue
}

func (s *StepStopInstance) Cleanup(_ multistep.StateBag) {}
