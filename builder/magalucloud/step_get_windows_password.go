// Copyright (c) Juliano Fernandes 2026
// SPDX-License-Identifier: MPL-2.0

package magalucloud

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/MagaluCloud/mgc-sdk-go/compute"
	"github.com/hashicorp/packer-plugin-sdk/communicator"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/hashicorp/packer-plugin-sdk/packer"
)

type StepGetWindowsPassword struct {
	Client *compute.VirtualMachineClient
	WinRM  *communicator.WinRM
}

func (s *StepGetWindowsPassword) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	id := state.Get("instance_id").(string)
	ui := state.Get("ui").(packer.Ui)
	ui.Say("Retrieving Windows password...")

	ticker := time.NewTicker(WaitInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			state.Put("error", ctx.Err())
			return multistep.ActionHalt
		case <-ticker.C:
			res, err := s.Client.Instances().GetFirstWindowsPassword(ctx, id)
			if err != nil && !strings.Contains(err.Error(), "windows_password_not_ready") {
				state.Put("error", fmt.Errorf("error retrieving windows password: %s", err))
				return multistep.ActionHalt
			}
			if err == nil {
				if s.WinRM.WinRMUser == "" {
					s.WinRM.WinRMUser = res.Instance.User
				}
				if s.WinRM.WinRMUser != res.Instance.User {
					state.Put("error", fmt.Errorf("winrm_username mismatch: %s != %s", s.WinRM.WinRMUser, res.Instance.User))
					return multistep.ActionHalt
				}

				s.WinRM.WinRMPassword = res.Instance.Password
				log.Printf("[DEBUG] Username: %s Password: %s", s.WinRM.WinRMUser, s.WinRM.WinRMPassword)

				return multistep.ActionContinue
			}
		}
	}
}

func (s *StepGetWindowsPassword) Cleanup(_ multistep.StateBag) {}
