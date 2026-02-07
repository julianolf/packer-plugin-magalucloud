package magalucloud

import (
	"context"
	"fmt"

	"github.com/MagaluCloud/mgc-sdk-go/compute"
	"github.com/MagaluCloud/mgc-sdk-go/helpers"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/hashicorp/packer-plugin-sdk/packer"
)

type StepCreateSnapshot struct{}

func (s *StepCreateSnapshot) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	ui := state.Get("ui").(packer.Ui)
	cfg := state.Get("config").(*Config)
	cli := state.Get("compute").(*compute.VirtualMachineClient)
	id := state.Get("instance_id").(string)

	ui.Say(fmt.Sprintf("Creating virtual machine snapshot: %s", id))

	req := compute.CreateSnapshotRequest{
		Name:     cfg.ImageName,
		Instance: compute.IDOrName{ID: helpers.StrPtr(id)},
	}

	id, err := cli.Snapshots().Create(ctx, req)
	if err != nil {
		state.Put("error", fmt.Errorf("Error creating snapshot: %s", err))
		return multistep.ActionHalt
	}

	state.Put("snapshot_id", id)
	return multistep.ActionContinue
}

func (s *StepCreateSnapshot) Cleanup(_ multistep.StateBag) {}
