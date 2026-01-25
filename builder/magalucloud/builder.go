// Copyright (c) Juliano Fernandes 2026
// SPDX-License-Identifier: MPL-2.0

//go:generate packer-sdc mapstructure-to-hcl2 -type Config

package magalucloud

import (
	"context"
	"fmt"
	"time"

	"github.com/MagaluCloud/mgc-sdk-go/client"
	"github.com/MagaluCloud/mgc-sdk-go/compute"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/packer-plugin-sdk/common"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/hashicorp/packer-plugin-sdk/multistep/commonsteps"
	"github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/hashicorp/packer-plugin-sdk/template/config"
)

const (
	BuilderId    = "julianolf.magalucloud"
	waitInterval = 30 * time.Second
)

type Region string

var Regions map[Region]client.MgcUrl = map[Region]client.MgcUrl{
	"br-se1": client.BrSe1,
	"br-ne1": client.BrNe1,
}

type Config struct {
	common.PackerConfig `mapstructure:",squash"`
	Token               string        `mapstructure:"token" required:"true"`
	Region              Region        `mapstructure:"region" required:"true"`
	SourceImage         string        `mapstructure:"source_image" required:"true"`
	MachineType         string        `mapstructure:"machine_type" required:"true"`
	SSHKey              string        `mapstructure:"ssh_key" required:"true"`
	ImageName           string        `mapstructure:"image_name" required:"false"`
	URL                 client.MgcUrl `mapstructure:"url" required:"false"`
}

type Builder struct {
	config Config
	runner multistep.Runner
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

	url, ok := Regions[b.config.Region]
	if !ok {
		return nil, nil, fmt.Errorf("Invalid region: %s", b.config.Region)
	}

	b.config.URL = url
	b.config.ImageName = fmt.Sprintf("packer-%s", time.Now().Format("20060102150405"))
	return nil, nil, nil
}

func (b *Builder) Run(ctx context.Context, ui packer.Ui, hook packer.Hook) (packer.Artifact, error) {
	cli := client.NewMgcClient(client.WithBaseURL(b.config.URL), client.WithAPIKey(b.config.Token))

	state := &multistep.BasicStateBag{}
	state.Put("config", &b.config)
	state.Put("hook", hook)
	state.Put("ui", ui)
	state.Put("compute", compute.New(cli))

	steps := []multistep.Step{
		&StepCreateInstance{},
		&StepWaitInstanceBoot{},
		&commonsteps.StepProvision{},
		&StepDeleteInstance{},
		&StepWaitInstanceTeardown{},
	}

	b.runner = commonsteps.NewRunner(steps, b.config.PackerConfig, ui)
	b.runner.Run(ctx, state)

	if err, ok := state.GetOk("error"); ok {
		return nil, err.(error)
	}

	artifact := &Artifact{StateData: map[string]any{}}
	return artifact, nil
}
