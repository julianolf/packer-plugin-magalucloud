// Copyright (c) Juliano Fernandes 2026
// SPDX-License-Identifier: MPL-2.0

package magalucloud

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

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
	sg := state.Get("security_group_id").(string)
	ui := state.Get("ui").(packer.Ui)
	ui.Sayf("Creating virtual machine instance from %s", s.Config.SourceImage)

	req := compute.CreateRequest{
		Name:             s.Config.ImageName,
		MachineType:      compute.IDOrName{Name: helpers.StrPtr(s.Config.MachineType)},
		Image:            compute.IDOrName{Name: helpers.StrPtr(s.Config.SourceImage)},
		AvailabilityZone: helpers.StrPtr(s.Config.AvailabilityZone),
		Network: &compute.CreateParametersNetwork{
			AssociatePublicIp: helpers.BoolPtr(true),
			Interface: &compute.CreateParametersNetworkInterface{
				SecurityGroups: &[]compute.CreateParametersNetworkInterfaceWithID{
					{ID: sg},
				},
			},
		},
	}
	if s.Config.Comm.Type == "ssh" {
		req.SshKeyName = helpers.StrPtr(s.Config.Comm.SSHTemporaryKeyPairName)
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

func (s *StepCreateInstance) Cleanup(state multistep.StateBag) {
	v, ok := state.GetOk("instance_id")
	if !ok {
		return
	}

	id := v.(string)
	ui := state.Get("ui").(packer.Ui)
	ui.Sayf("Deleting virtual machine instance %s", id)

	err := s.Client.Instances().Delete(context.Background(), id, true)
	if err != nil {
		ui.Errorf("Error deleting virtual machine: %s", err)
		return
	}

	ui.Sayf("Waiting for virtual machine instance %s teardown", id)

	ticker := time.NewTicker(WaitInterval)
	defer ticker.Stop()

	timeout := time.NewTimer(TimeoutInterval)
	defer timeout.Stop()

	for {
		select {
		case <-timeout.C:
			ui.Errorf("Delete virtual machine %s timed out", id)
			return
		case <-ticker.C:
			instance, err := s.Client.Instances().Get(context.Background(), id, []compute.InstanceExpand{})
			if err != nil && strings.Contains(err.Error(), "404") {
				return
			}
			if err != nil {
				ui.Errorf("Error querying virtual machine: %s", err)
				return
			}
			if instance.State == "error" {
				ui.Errorf("Virtual machine state error: %s", instance.Error.Message)
				return
			}
			if strings.Contains(instance.Status, "error") {
				ui.Errorf("Virtual machine status error: %s", instance.Error.Message)
				return
			}
		}
	}
}
