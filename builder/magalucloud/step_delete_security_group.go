// Copyright (c) Juliano Fernandes 2026
// SPDX-License-Identifier: MPL-2.0

package magalucloud

import (
	"context"
	"fmt"

	"github.com/MagaluCloud/mgc-sdk-go/network"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/hashicorp/packer-plugin-sdk/packer"
)

type StepDeleteSecurityGroup struct {
	Client *network.NetworkClient
}

func (s *StepDeleteSecurityGroup) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	id := state.Get("security_group_id").(string)
	ui := state.Get("ui").(packer.Ui)
	ui.Sayf("Deleting security group %s", id)

	err := s.Client.SecurityGroups().Delete(ctx, id)
	if err != nil {
		state.Put("error", fmt.Errorf("error deleting security group: %s", err))
		return multistep.ActionHalt
	}

	return multistep.ActionContinue
}

func (s *StepDeleteSecurityGroup) Cleanup(_ multistep.StateBag) {}
