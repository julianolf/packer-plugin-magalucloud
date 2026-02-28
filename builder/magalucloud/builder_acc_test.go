// Copyright (c) Juliano Fernandes 2026
// SPDX-License-Identifier: MPL-2.0

package magalucloud

import (
	_ "embed"
	"fmt"
	"strings"

	"os"
	"os/exec"
	"testing"

	"github.com/hashicorp/packer-plugin-sdk/acctest"
)

//go:embed test-fixtures/template.pkr.hcl
var testBuilderHCL2Basic string

// Run with: PACKER_ACC=1 go test -count 1 -v ./builder/magalucloud/builder_acc_test.go  -timeout=120m
func TestAccMagaluCloudBuilder(t *testing.T) {
	testCase := &acctest.PluginTestCase{
		Name: "magalucloud_builder_test",
		Setup: func() error {
			return nil
		},
		Teardown: func() error {
			return nil
		},
		Template: testBuilderHCL2Basic,
		Type:     "magalucloud",
		Check: func(buildCommand *exec.Cmd, logfile string) error {
			if buildCommand.ProcessState != nil {
				if buildCommand.ProcessState.ExitCode() != 0 {
					return fmt.Errorf("Bad exit code. Logfile: %s", logfile)
				}
			}

			data, err := os.ReadFile(logfile)
			if err != nil {
				return fmt.Errorf("Unable to read %s", logfile)
			}

			logs := string(data)
			expectedLogs := []string{
				"test.magalucloud.basic: Creating temporary ED25519 SSH key",
				"test.magalucloud.basic: Uploading SSH key",
				"test.magalucloud.basic: Creating virtual machine instance from cloud-ubuntu-22.04 LTS",
				"test.magalucloud.basic: Waiting for virtual machine instance",
				"test.magalucloud.basic: Using SSH communicator to connect",
				"test.magalucloud.basic: Waiting for SSH to become available",
				"test.magalucloud.basic: Connected to SSH",
				"test.magalucloud.basic: Provisioning with shell script",
				"test.magalucloud.basic: This is a test",
				"test.magalucloud.basic: Stopping virtual machine instance",
				"test.magalucloud.basic: Creating a snapshot of the virtual machine instance",
				"test.magalucloud.basic: Waiting for snapshot",
				"test.magalucloud.basic: Deleting virtual machine instance",
				"test.magalucloud.basic: Deleting SSH key",
			}

			for _, expected := range expectedLogs {
				if !strings.Contains(logs, expected) {
					t.Fatalf("logs doesn't contain expected value %s", expected)
				}
			}

			return nil
		},
	}
	acctest.TestPlugin(t, testCase)
}
