// Copyright (c) Juliano Fernandes 2026
// SPDX-License-Identifier: MPL-2.0

package magalucloud

import (
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/hashicorp/packer-plugin-sdk/acctest"
)

//go:embed test-fixtures/template.pkr.hcl
var testBuilderHCL2 string

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
		Template: testBuilderHCL2,
		Type:     "magalucloud",
		Check: func(buildCommand *exec.Cmd, logfile string) error {
			if buildCommand.ProcessState != nil {
				if buildCommand.ProcessState.ExitCode() != 0 {
					return fmt.Errorf("bad exit code. logfile: %s", logfile)
				}
			}

			data, err := os.ReadFile(logfile)
			if err != nil {
				return fmt.Errorf("unable to read %s", logfile)
			}

			logs := string(data)
			expectedLogs := []string{
				"test.magalucloud.test: Creating temporary ED25519 SSH key",
				"test.magalucloud.test: Uploading SSH key",
				"test.magalucloud.test: Creating security group",
				"test.magalucloud.test: Creating virtual machine instance from cloud-debian-12 LTS",
				"test.magalucloud.test: Waiting for virtual machine instance",
				"test.magalucloud.test: Using SSH communicator to connect",
				"test.magalucloud.test: Waiting for SSH to become available",
				"test.magalucloud.test: Connected to SSH",
				"test.magalucloud.test: Provisioning with shell script",
				"test.magalucloud.test: This is a test",
				"test.magalucloud.test: Stopping virtual machine instance",
				"test.magalucloud.test: Creating a snapshot of the virtual machine instance",
				"test.magalucloud.test: Waiting for snapshot",
				"test.magalucloud.test: Deleting virtual machine instance",
				"test.magalucloud.test: Deleting security group",
				"test.magalucloud.test: Deleting SSH key",
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
