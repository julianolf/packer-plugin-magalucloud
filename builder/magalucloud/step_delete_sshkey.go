// Copyright (c) Juliano Fernandes 2026
// SPDX-License-Identifier: MPL-2.0

package magalucloud

import (
	"context"
	"fmt"

	"github.com/MagaluCloud/mgc-sdk-go/sshkeys"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/hashicorp/packer-plugin-sdk/packer"
)

type StepDeleteSSHKey struct{}

func (s *StepDeleteSSHKey) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	ui := state.Get("ui").(packer.Ui)
	cfg := state.Get("config").(*Config)
	cli := state.Get("sshkeys").(*sshkeys.SSHKeyClient)
	id := state.Get("sshkey_id").(string)

	ui.Say(fmt.Sprintf("Deleting SSH key: %s", cfg.Comm.SSH.SSHTemporaryKeyPairName))

	_, err := cli.Keys().Delete(ctx, id)
	if err != nil {
		state.Put("error", fmt.Errorf("Error deleting SSH key: %s", err))
		return multistep.ActionHalt
	}

	return multistep.ActionContinue
}

func (s *StepDeleteSSHKey) Cleanup(_ multistep.StateBag) {}
