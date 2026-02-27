// Copyright (c) Juliano Fernandes 2026
// SPDX-License-Identifier: MPL-2.0

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

type StepWaitInstanceBoot struct {
	Client *compute.VirtualMachineClient
}

func (s *StepWaitInstanceBoot) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	id := state.Get("instance_id").(string)
	ui := state.Get("ui").(packer.Ui)
	ui.Sayf("Waiting for virtual machine instance %s boot", id)

	ticker := time.NewTicker(WaitInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			state.Put("error", ctx.Err())
			return multistep.ActionHalt
		case <-ticker.C:
			instance, err := s.Client.Instances().Get(ctx, id, []compute.InstanceExpand{compute.InstanceNetworkExpand})
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
			if instance.State == "running" && instance.Status == "completed" {
				ip := (*instance.Network.Interfaces)[0].AssociatedPublicIpv4
				if ip == nil {
					continue
				}
				state.Put("instance_ip", *ip)
				return multistep.ActionContinue
			}
		}
	}
}

func (s *StepWaitInstanceBoot) Cleanup(_ multistep.StateBag) {}
