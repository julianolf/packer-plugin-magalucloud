// Copyright (c) Juliano Fernandes 2026
// SPDX-License-Identifier: MPL-2.0

package magalucloud

import (
	"context"
	"fmt"

	"github.com/MagaluCloud/mgc-sdk-go/sshkeys"
	"github.com/hashicorp/packer-plugin-sdk/communicator"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/hashicorp/packer-plugin-sdk/packer"
)

type StepDeleteSSHKey struct {
	Client *sshkeys.SSHKeyClient
	SSH    *communicator.SSH
}

func (s *StepDeleteSSHKey) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	ui := state.Get("ui").(packer.Ui)
	ui.Sayf("Deleting SSH key %s", s.SSH.SSHTemporaryKeyPairName)

	id := state.Get("sshkey_id").(string)
	_, err := s.Client.Keys().Delete(ctx, id)
	if err != nil {
		state.Put("error", fmt.Errorf("error deleting ssh key: %s", err))
		return multistep.ActionHalt
	}

	return multistep.ActionContinue
}

func (s *StepDeleteSSHKey) Cleanup(_ multistep.StateBag) {}
