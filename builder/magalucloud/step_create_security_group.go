// Copyright (c) Juliano Fernandes 2026
// SPDX-License-Identifier: MPL-2.0

package magalucloud

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/MagaluCloud/mgc-sdk-go/helpers"
	"github.com/MagaluCloud/mgc-sdk-go/network"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/hashicorp/packer-plugin-sdk/packer"
	"golang.org/x/sync/errgroup"
)

type StepCreateSecurityGroup struct {
	Client *network.NetworkClient
	Config *Config
}

func (s *StepCreateSecurityGroup) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	ui := state.Get("ui").(packer.Ui)
	ui.Sayf("Creating security group %s", s.Config.uname)

	req := network.SecurityGroupCreateRequest{
		Name:        s.Config.uname,
		Description: helpers.StrPtr(fmt.Sprintf("Security group for Packer build: %s", s.Config.PackerBuildName)),
	}

	data, _ := json.Marshal(req)
	log.Printf("[DEBUG] Create security group data %s", data)

	id, err := s.Client.SecurityGroups().Create(ctx, req)
	if err != nil {
		state.Put("error", fmt.Errorf("error creating security group: %s", err))
		return multistep.ActionHalt
	}
	log.Printf("[DEBUG] Security group ID %s", id)

	state.Put("security_group_id", id)

	rules := make([]network.RuleCreateRequest, 0, 2)
	switch s.Config.Comm.Type {
	case "ssh":
		rules = append(
			rules,
			network.RuleCreateRequest{
				Direction:      helpers.StrPtr("ingress"),
				PortRangeMin:   helpers.IntPtr(22),
				PortRangeMax:   helpers.IntPtr(22),
				Protocol:       helpers.StrPtr("tcp"),
				RemoteIPPrefix: helpers.StrPtr("0.0.0.0/0"),
				EtherType:      *helpers.StrPtr("IPv4"),
			},
			network.RuleCreateRequest{
				Direction:      helpers.StrPtr("ingress"),
				PortRangeMin:   helpers.IntPtr(22),
				PortRangeMax:   helpers.IntPtr(22),
				Protocol:       helpers.StrPtr("tcp"),
				RemoteIPPrefix: helpers.StrPtr("::/0"),
				EtherType:      *helpers.StrPtr("IPv6"),
			},
		)
	case "winrm":
		rules = append(
			rules,
			network.RuleCreateRequest{
				Direction:      helpers.StrPtr("ingress"),
				PortRangeMin:   helpers.IntPtr(5985),
				PortRangeMax:   helpers.IntPtr(5986),
				Protocol:       helpers.StrPtr("tcp"),
				RemoteIPPrefix: helpers.StrPtr("0.0.0.0/0"),
				EtherType:      *helpers.StrPtr("IPv4"),
			},
			network.RuleCreateRequest{
				Direction:      helpers.StrPtr("ingress"),
				PortRangeMin:   helpers.IntPtr(5985),
				PortRangeMax:   helpers.IntPtr(5986),
				Protocol:       helpers.StrPtr("tcp"),
				RemoteIPPrefix: helpers.StrPtr("::/0"),
				EtherType:      *helpers.StrPtr("IPv6"),
			},
		)
	}

	data, _ = json.Marshal(rules)
	log.Printf("[DEBUG] Create security group rules data %s", data)

	g, ctx := errgroup.WithContext(ctx)
	for _, rule := range rules {
		r := rule
		g.Go(func() error {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}

			_, err := s.Client.Rules().Create(ctx, id, r)
			return err
		})
	}
	if err := g.Wait(); err != nil {
		state.Put("error", fmt.Errorf("error creating security group rules: %s", err))
		return multistep.ActionHalt
	}

	return multistep.ActionContinue
}

func (s *StepCreateSecurityGroup) Cleanup(_ multistep.StateBag) {}
