// Copyright (c) Juliano Fernandes 2026
// SPDX-License-Identifier: MPL-2.0

package magalucloud

import (
	_ "embed"
	"fmt"
	"io"

	"os"
	"os/exec"
	"regexp"
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

			logs, err := os.Open(logfile)
			if err != nil {
				return fmt.Errorf("Unable find %s", logfile)
			}
			defer logs.Close()

			logsBytes, err := io.ReadAll(logs)
			if err != nil {
				return fmt.Errorf("Unable to read %s", logfile)
			}
			logsString := string(logsBytes)

			buildGeneratedDataLog := "test.magalucloud.basic: Creating virtual machine instance from: cloud-ubuntu-22.04 LTS"
			if matched, _ := regexp.MatchString(buildGeneratedDataLog+".*", logsString); !matched {
				t.Fatalf("logs doesn't contain expected foo value %q", logsString)
			}
			return nil
		},
	}
	acctest.TestPlugin(t, testCase)
}
