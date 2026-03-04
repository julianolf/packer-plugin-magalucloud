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
var testImporterHCL2 string

// Run with: PACKER_ACC=1 go test -count 1 -v ./post-processor/magalucloud/import_acc_test.go  -timeout=120m
func TestAccMagaluCloudImporter(t *testing.T) {
	testCase := &acctest.PluginTestCase{
		Name: "magalucloud_importer_test",
		Setup: func() error {
			return nil
		},
		Teardown: func() error {
			return nil
		},
		Template: testImporterHCL2,
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
				"test.file.test (magalucloud-import): Uploading",
				"test.file.test (magalucloud-import): Generating presigned URL",
				"test.file.test (magalucloud-import): Importing image",
				"test.file.test (magalucloud-import): Finished importing image",
				"test.file.test (magalucloud-import): Deleting",
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
