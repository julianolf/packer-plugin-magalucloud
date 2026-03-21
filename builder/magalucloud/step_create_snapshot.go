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

type StepCreateSnapshot struct {
	Client *compute.VirtualMachineClient
	Config *Config
}

func (s *StepCreateSnapshot) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	id := state.Get("instance_id").(string)
	ui := state.Get("ui").(packer.Ui)
	ui.Sayf("Creating a snapshot of the virtual machine instance %s", id)

	req := compute.CreateSnapshotRequest{
		Name:     s.Config.ImageName,
		Instance: compute.IDOrName{ID: helpers.StrPtr(id)},
	}

	data, _ := json.Marshal(req)
	log.Printf("[DEBUG] Create snapshot data %s", data)

	id, err := s.Client.Snapshots().Create(ctx, req)
	if err != nil {
		state.Put("error", fmt.Errorf("error creating snapshot: %s", err))
		return multistep.ActionHalt
	}

	state.Put("snapshot_id", id)

	ui.Sayf("Waiting for snapshot %s creation", id)

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
			state.Put("error", fmt.Errorf("create snapshot %s timed out after: %s", id, s.Config.WaitTimeout))
			return multistep.ActionHalt
		case <-ticker.C:
			snapshot, err := s.Client.Snapshots().Get(ctx, id, []compute.SnapshotExpand{})
			if err != nil {
				state.Put("error", fmt.Errorf("error querying snapshot: %s", err))
				return multistep.ActionHalt
			}
			if strings.Contains(snapshot.Status, "error") {
				state.Put("error", fmt.Errorf("snapshot status error: %s", snapshot.Status))
				return multistep.ActionHalt
			}
			if snapshot.State == "available" && snapshot.Status == "completed" {
				return multistep.ActionContinue
			}
		}
	}
}

func (s *StepCreateSnapshot) Cleanup(_ multistep.StateBag) {}
