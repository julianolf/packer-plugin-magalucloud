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

type StepStopInstance struct {
	Client *compute.VirtualMachineClient
}

func (s *StepStopInstance) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	id := state.Get("instance_id").(string)
	ui := state.Get("ui").(packer.Ui)
	ui.Sayf("Stopping virtual machine instance %s", id)

	err := s.Client.Instances().Stop(ctx, id)
	if err != nil {
		state.Put("error", fmt.Errorf("Error stopping virtual machine: %s", err))
		return multistep.ActionHalt
	}

	return multistep.ActionContinue
}

func (s *StepStopInstance) Cleanup(_ multistep.StateBag) {}
