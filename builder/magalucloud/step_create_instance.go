// Copyright (c) Juliano Fernandes 2026
// SPDX-License-Identifier: MPL-2.0

package magalucloud

import (
	"context"
	"fmt"

	"github.com/MagaluCloud/mgc-sdk-go/compute"
	"github.com/MagaluCloud/mgc-sdk-go/helpers"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/hashicorp/packer-plugin-sdk/packer"
)

type StepCreateInstance struct{}

func (s *StepCreateInstance) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	ui := state.Get("ui").(packer.Ui)
	cfg := state.Get("config").(*Config)
	cli := state.Get("compute").(*compute.VirtualMachineClient)

	ui.Say(fmt.Sprintf("Creating virtual machine instance from: %s", cfg.SourceImage))

	req := compute.CreateRequest{
		Name:        cfg.ImageName,
		MachineType: compute.IDOrName{Name: helpers.StrPtr(cfg.MachineType)},
		Image:       compute.IDOrName{Name: helpers.StrPtr(cfg.SourceImage)},
		Network:     &compute.CreateParametersNetwork{AssociatePublicIp: helpers.BoolPtr(true)},
		SshKeyName:  helpers.StrPtr(cfg.Comm.SSH.SSHTemporaryKeyPairName),
	}

	id, err := cli.Instances().Create(ctx, req)
	if err != nil {
		state.Put("error", fmt.Errorf("Error creating virtual machine: %s", err))
		return multistep.ActionHalt
	}

	state.Put("instance_id", id)
	return multistep.ActionContinue
}

func (s *StepCreateInstance) Cleanup(_ multistep.StateBag) {}
