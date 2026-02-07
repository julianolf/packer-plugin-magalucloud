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

type StepWaitInstanceStop struct{}

func (s *StepWaitInstanceStop) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	ui := state.Get("ui").(packer.Ui)
	cli := state.Get("compute").(*compute.VirtualMachineClient)
	id := state.Get("instance_id").(string)

	ui.Say(fmt.Sprintf("Waiting for virtual machine instance with ID: %s", id))

	ticker := time.NewTicker(waitInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			state.Put("error", ctx.Err())
			return multistep.ActionHalt
		case <-ticker.C:
			instance, err := cli.Instances().Get(ctx, id, []compute.InstanceExpand{})
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
			if instance.State == "stopped" && instance.Status == "completed" {
				return multistep.ActionContinue
			}
		}
	}
}

func (s *StepWaitInstanceStop) Cleanup(_ multistep.StateBag) {}
