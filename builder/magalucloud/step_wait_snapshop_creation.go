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

type StepWaitSnapshotCreation struct{}

func (s *StepWaitSnapshotCreation) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	ui := state.Get("ui").(packer.Ui)
	cli := state.Get("compute").(*compute.VirtualMachineClient)
	id := state.Get("snapshot_id").(string)

	ui.Say(fmt.Sprintf("Waiting for snapshot with ID: %s", id))

	ticker := time.NewTicker(waitInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			state.Put("error", ctx.Err())
			return multistep.ActionHalt
		case <-ticker.C:
			snapshot, err := cli.Snapshots().Get(ctx, id, []compute.SnapshotExpand{})
			if err != nil {
				state.Put("error", fmt.Errorf("Error querying snapshot: %s", err))
				return multistep.ActionHalt
			}
			if strings.Contains(snapshot.Status, "error") {
				state.Put("error", fmt.Errorf("Snapshot status error: %s", snapshot.Status))
				return multistep.ActionHalt
			}
			if snapshot.State == "available" && snapshot.Status == "completed" {
				return multistep.ActionContinue
			}
		}
	}
}

func (s *StepWaitSnapshotCreation) Cleanup(_ multistep.StateBag) {}
