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

type StepUploadSSHKey struct{}

func (s *StepUploadSSHKey) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	ui := state.Get("ui").(packer.Ui)
	cfg := state.Get("config").(*Config)
	cli := state.Get("sshkeys").(*sshkeys.SSHKeyClient)

	ui.Say(fmt.Sprintf("Uploading SSH key: %s", cfg.Comm.SSH.SSHTemporaryKeyPairName))

	req := sshkeys.CreateSSHKeyRequest{
		Name: cfg.Comm.SSH.SSHTemporaryKeyPairName,
		Key:  string(cfg.Comm.SSH.SSHPublicKey),
	}

	key, err := cli.Keys().Create(ctx, req)
	if err != nil {
		state.Put("error", fmt.Errorf("Error uploading SSH key: %s", err))
		return multistep.ActionHalt
	}

	state.Put("sshkey_id", key.ID)
	return multistep.ActionContinue
}

func (s *StepUploadSSHKey) Cleanup(_ multistep.StateBag) {}
