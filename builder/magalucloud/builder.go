// Copyright (c) Juliano Fernandes 2026
// SPDX-License-Identifier: MPL-2.0

//go:generate packer-sdc mapstructure-to-hcl2 -type Config

package magalucloud

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/MagaluCloud/mgc-sdk-go/client"
	"github.com/MagaluCloud/mgc-sdk-go/compute"
	"github.com/MagaluCloud/mgc-sdk-go/network"
	"github.com/MagaluCloud/mgc-sdk-go/sshkeys"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/packer-plugin-sdk/common"
	"github.com/hashicorp/packer-plugin-sdk/communicator"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/hashicorp/packer-plugin-sdk/multistep/commonsteps"
	"github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/hashicorp/packer-plugin-sdk/template/config"
	"github.com/hashicorp/packer-plugin-sdk/template/interpolate"
	"github.com/hashicorp/packer-plugin-sdk/uuid"
)

const (
	BuilderId      = "julianolf.magalucloud"
	WaitInterval   = 10 * time.Second
	DefaultTimeout = 10 * time.Minute
)

type Region string

var Regions map[Region]client.MgcUrl = map[Region]client.MgcUrl{
	"br-se1": client.BrSe1,
	"br-ne1": client.BrNe1,
}

type Config struct {
	common.PackerConfig `mapstructure:",squash"`
	Comm                communicator.Config `mapstructure:",squash"`
	APIKey              string              `mapstructure:"api_key"`
	Region              Region              `mapstructure:"region"`
	AvailabilityZone    string              `mapstructure:"availability_zone"`
	SourceImage         string              `mapstructure:"source_image"`
	MachineType         string              `mapstructure:"machine_type"`
	ImageName           string              `mapstructure:"image_name"`
	URL                 client.MgcUrl       `mapstructure:"url"`
	WaitTimeout         time.Duration       `mapstructure:"wait_timeout"`

	uname string
	ctx   interpolate.Context
}

func (c *Config) WinRMConfigFunc(state multistep.StateBag) (*communicator.WinRMConfig, error) {
	return &communicator.WinRMConfig{
		Username: c.Comm.WinRMUser,
		Password: c.Comm.WinRMPassword,
	}, nil
}

type Builder struct {
	config  Config
	runner  multistep.Runner
	sshkeys *sshkeys.SSHKeyClient
	compute *compute.VirtualMachineClient
	network *network.NetworkClient
}

func (b *Builder) ConfigSpec() hcldec.ObjectSpec {
	return b.config.FlatMapstructure().HCL2Spec()
}

func (b *Builder) Prepare(raws ...any) (generatedVars []string, warnings []string, err error) {
	err = config.Decode(
		&b.config,
		&config.DecodeOpts{
			PluginType:  BuilderId,
			Interpolate: true,
		},
		raws...,
	)
	if err != nil {
		return nil, nil, err
	}

	if err := b.config.Comm.Prepare(&b.config.ctx); len(err) > 0 {
		return nil, nil, errors.Join(err...)
	}

	url, ok := Regions[b.config.Region]
	if !ok {
		return nil, nil, fmt.Errorf("invalid region: %s", b.config.Region)
	}
	b.config.URL = url

	if b.config.AvailabilityZone == "" {
		b.config.AvailabilityZone = fmt.Sprintf("%s-a", b.config.Region)
	}
	if !strings.HasPrefix(b.config.AvailabilityZone, string(b.config.Region)) {
		return nil, nil, fmt.Errorf("invalid availability zone: %s", b.config.AvailabilityZone)
	}

	if b.config.WaitTimeout == time.Duration(0) {
		b.config.WaitTimeout = DefaultTimeout
	}

	b.config.uname = fmt.Sprintf("packer-%s", uuid.TimeOrderedUUID())

	if b.config.ImageName == "" {
		b.config.ImageName = b.config.uname
	}

	b.config.Comm.SSHTemporaryKeyPairName = b.config.uname
	b.config.Comm.SSHTemporaryKeyPairType = "ed25519"

	b.sshkeys = sshkeys.New(client.NewMgcClient(client.WithAPIKey(b.config.APIKey)))
	b.compute = compute.New(client.NewMgcClient(client.WithAPIKey(b.config.APIKey), client.WithBaseURL(b.config.URL)))
	b.network = network.New(client.NewMgcClient(client.WithAPIKey(b.config.APIKey), client.WithBaseURL(b.config.URL)))

	return nil, nil, nil
}

func (b *Builder) Run(ctx context.Context, ui packer.Ui, hook packer.Hook) (packer.Artifact, error) {
	state := &multistep.BasicStateBag{}
	state.Put("hook", hook)
	state.Put("ui", ui)

	steps := []multistep.Step{
		multistep.If(
			b.config.Comm.Type == "ssh",
			&communicator.StepSSHKeyGen{
				CommConf:            &b.config.Comm,
				SSHTemporaryKeyPair: b.config.Comm.SSHTemporaryKeyPair,
			},
		),
		multistep.If(
			b.config.PackerDebug && b.config.Comm.Type == "ssh",
			&communicator.StepDumpSSHKey{
				Path: fmt.Sprintf("mgc_%s.pem", b.config.PackerBuildName),
				SSH:  &b.config.Comm.SSH,
			},
		),
		multistep.If(
			b.config.Comm.Type == "ssh",
			&StepUploadSSHKey{
				Client: b.sshkeys,
				SSH:    &b.config.Comm.SSH,
			},
		),
		&StepCreateSecurityGroup{
			Client: b.network,
			Config: &b.config,
		},
		&StepCreateInstance{
			Client: b.compute,
			Config: &b.config,
		},
		multistep.If(
			b.config.Comm.Type == "winrm",
			&StepGetWindowsPassword{
				Client: b.compute,
				Config: &b.config,
			},
		),
		&communicator.StepConnect{
			Config:      &b.config.Comm,
			Host:        communicator.CommHost(b.config.Comm.Host(), "instance_ip"),
			SSHConfig:   b.config.Comm.SSHConfigFunc(),
			WinRMConfig: b.config.WinRMConfigFunc,
		},
		&commonsteps.StepProvision{},
		&StepStopInstance{
			Client: b.compute,
			Config: &b.config,
		},
		&StepCreateSnapshot{
			Client: b.compute,
			Config: &b.config,
		},
	}

	b.runner = commonsteps.NewRunner(steps, b.config.PackerConfig, ui)
	b.runner.Run(ctx, state)

	if err, ok := state.GetOk("error"); ok {
		return nil, err.(error)
	}

	snapshot_id := state.Get("snapshot_id").(string)
	artifact := &Artifact{
		ID:        snapshot_id,
		Region:    b.config.Region,
		StateData: map[string]any{},
	}
	return artifact, nil
}
