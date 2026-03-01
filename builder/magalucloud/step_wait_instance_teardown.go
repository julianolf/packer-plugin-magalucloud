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

type StepWaitInstanceTeardown struct {
	Client *compute.VirtualMachineClient
}

func (s *StepWaitInstanceTeardown) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	id := state.Get("instance_id").(string)
	ui := state.Get("ui").(packer.Ui)
	ui.Sayf("Waiting for virtual machine instance %s teardown", id)

	ticker := time.NewTicker(WaitInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			state.Put("error", ctx.Err())
			return multistep.ActionHalt
		case <-ticker.C:
			instance, err := s.Client.Instances().Get(ctx, id, []compute.InstanceExpand{})
			if err != nil && strings.Contains(err.Error(), "404") {
				return multistep.ActionContinue
			}
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
		}
	}
}

func (s *StepWaitInstanceTeardown) Cleanup(_ multistep.StateBag) {}
