// Copyright (c) Juliano Fernandes 2026
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"fmt"
	"os"

	"github.com/julianolf/packer-plugin-magalucloud/builder/magalucloud"
	magalucloudData "github.com/julianolf/packer-plugin-magalucloud/datasource/magalucloud"
	magalucloudPP "github.com/julianolf/packer-plugin-magalucloud/post-processor/magalucloud"
	magalucloudProv "github.com/julianolf/packer-plugin-magalucloud/provisioner/magalucloud"
	magalucloudVersion "github.com/julianolf/packer-plugin-magalucloud/version"

	"github.com/hashicorp/packer-plugin-sdk/plugin"
)

func main() {
	pps := plugin.NewSet()
	pps.RegisterBuilder("my-builder", new(magalucloud.Builder))
	pps.RegisterProvisioner("my-provisioner", new(magalucloudProv.Provisioner))
	pps.RegisterPostProcessor("my-post-processor", new(magalucloudPP.PostProcessor))
	pps.RegisterDatasource("my-datasource", new(magalucloudData.Datasource))
	pps.SetVersion(magalucloudVersion.PluginVersion)
	err := pps.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
