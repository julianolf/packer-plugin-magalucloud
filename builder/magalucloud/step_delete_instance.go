package magalucloud

import (
	"context"
	"fmt"

	"github.com/MagaluCloud/mgc-sdk-go/compute"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/hashicorp/packer-plugin-sdk/packer"
)

type StepDeleteInstance struct{}

func (s *StepDeleteInstance) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	ui := state.Get("ui").(packer.Ui)
	cli := state.Get("compute").(*compute.VirtualMachineClient)
	id := state.Get("instance_id").(string)

	ui.Say(fmt.Sprintf("Deleting virtual machine instance with ID: %s", id))

	err := cli.Instances().Delete(ctx, id, true)
	if err != nil {
		state.Put("error", fmt.Errorf("Error deleting virtual machine: %s", err))
		return multistep.ActionHalt
	}

	return multistep.ActionContinue
}

func (s *StepDeleteInstance) Cleanup(_ multistep.StateBag) {}
