// Copyright (c) Juliano Fernandes 2026
// SPDX-License-Identifier: MPL-2.0

package magalucloud

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/MagaluCloud/mgc-sdk-go/compute"
	"github.com/MagaluCloud/mgc-sdk-go/helpers"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/hashicorp/packer-plugin-sdk/packer"
)

type StepCreateInstance struct {
	Client *compute.VirtualMachineClient
	Config *Config
}

func (s *StepCreateInstance) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	ui := state.Get("ui").(packer.Ui)
	ui.Sayf("Creating virtual machine instance from %s", s.Config.SourceImage)

	req := compute.CreateRequest{
		Name:             s.Config.ImageName,
		MachineType:      compute.IDOrName{Name: helpers.StrPtr(s.Config.MachineType)},
		Image:            compute.IDOrName{Name: helpers.StrPtr(s.Config.SourceImage)},
		Network:          &compute.CreateParametersNetwork{AssociatePublicIp: helpers.BoolPtr(true)},
		SshKeyName:       helpers.StrPtr(s.Config.Comm.SSHTemporaryKeyPairName),
		AvailabilityZone: helpers.StrPtr(s.Config.AvailabilityZone),
	}

	data, _ := json.Marshal(req)
	log.Printf("[DEBUG] Create instance data %s", data)

	id, err := s.Client.Instances().Create(ctx, req)
	if err != nil {
		state.Put("error", fmt.Errorf("error creating virtual machine: %s", err))
		return multistep.ActionHalt
	}

	state.Put("instance_id", id)
	return multistep.ActionContinue
}

func (s *StepCreateInstance) Cleanup(_ multistep.StateBag) {}
