// Copyright (c) Juliano Fernandes 2026
// SPDX-License-Identifier: MPL-2.0

package magalucloud

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/MagaluCloud/mgc-sdk-go/sshkeys"
	"github.com/hashicorp/packer-plugin-sdk/communicator"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/hashicorp/packer-plugin-sdk/packer"
)

type StepUploadSSHKey struct {
	Client *sshkeys.SSHKeyClient
	SSH    *communicator.SSH
}

func (s *StepUploadSSHKey) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	ui := state.Get("ui").(packer.Ui)
	ui.Sayf("Uploading SSH key %s", s.SSH.SSHTemporaryKeyPairName)

	req := sshkeys.CreateSSHKeyRequest{
		Name: s.SSH.SSHTemporaryKeyPairName,
		Key:  string(s.SSH.SSHPublicKey),
	}

	data, _ := json.Marshal(req)
	log.Printf("[DEBUG] Create SSH key data %s", data)

	key, err := s.Client.Keys().Create(ctx, req)
	if err != nil {
		state.Put("error", fmt.Errorf("error uploading ssh key: %s", err))
		return multistep.ActionHalt
	}

	state.Put("sshkey_id", key.ID)
	return multistep.ActionContinue
}

func (s *StepUploadSSHKey) Cleanup(_ multistep.StateBag) {}
