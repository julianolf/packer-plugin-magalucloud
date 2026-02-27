// Copyright (c) Juliano Fernandes 2026
// SPDX-License-Identifier: MPL-2.0

package magalucloud

import (
	"context"
	"fmt"

	"github.com/MagaluCloud/mgc-sdk-go/compute"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/hashicorp/packer-plugin-sdk/packer"
)

type StepDeleteInstance struct {
	Client *compute.VirtualMachineClient
}

func (s *StepDeleteInstance) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	id := state.Get("instance_id").(string)
	ui := state.Get("ui").(packer.Ui)
	ui.Sayf("Deleting virtual machine instance %s", id)

	err := s.Client.Instances().Delete(ctx, id, true)
	if err != nil {
		state.Put("error", fmt.Errorf("Error deleting virtual machine: %s", err))
		return multistep.ActionHalt
	}

	return multistep.ActionContinue
}

func (s *StepDeleteInstance) Cleanup(_ multistep.StateBag) {}
