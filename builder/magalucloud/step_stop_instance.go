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

type StepStopInstance struct {
	Client *compute.VirtualMachineClient
	Config *Config
}

func (s *StepStopInstance) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	id := state.Get("instance_id").(string)
	ui := state.Get("ui").(packer.Ui)
	ui.Sayf("Stopping virtual machine instance %s", id)

	err := s.Client.Instances().Stop(ctx, id)
	if err != nil {
		state.Put("error", fmt.Errorf("error stopping virtual machine: %s", err))
		return multistep.ActionHalt
	}

	ui.Sayf("Waiting for virtual machine instance %s stop", id)

	ticker := time.NewTicker(WaitInterval)
	defer ticker.Stop()

	timeout := time.NewTimer(s.Config.WaitTimeout)
	defer timeout.Stop()

	for {
		select {
		case <-ctx.Done():
			state.Put("error", ctx.Err())
			return multistep.ActionHalt
		case <-timeout.C:
			state.Put("error", fmt.Errorf("stop virtual machine %s timed out after: %s", id, s.Config.WaitTimeout))
			return multistep.ActionHalt
		case <-ticker.C:
			instance, err := s.Client.Instances().Get(ctx, id, []compute.InstanceExpand{})
			if err != nil {
				state.Put("error", fmt.Errorf("error querying virtual machine: %s", err))
				return multistep.ActionHalt
			}
			if instance.State == "error" {
				state.Put("error", fmt.Errorf("virtual machine state error: %s", instance.Error.Message))
				return multistep.ActionHalt
			}
			if strings.Contains(instance.Status, "error") {
				state.Put("error", fmt.Errorf("virtual machine status error: %s", instance.Error.Message))
				return multistep.ActionHalt
			}
			if instance.State == "stopped" && instance.Status == "completed" {
				return multistep.ActionContinue
			}
		}
	}
}

func (s *StepStopInstance) Cleanup(_ multistep.StateBag) {}
