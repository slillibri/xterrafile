package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/rendon/testcli"
	"github.com/stretchr/testify/assert"
)

var xTerrafileBinaryPath string
var workingDirectory string

func init() {
	var err error
	workingDirectory, err = os.Getwd()
	if err != nil {
		panic(err)
	}
	xTerrafileBinaryPath = workingDirectory + "/xterrafile"
}
func TestTerraformWithTerrafilePath(t *testing.T) {
	folder, back := setup(t)
	defer back()

	testcli.Run(xTerrafileBinaryPath, "-f", fmt.Sprint(folder, "/Terrafile.test"), "install")

	if !testcli.Success() {
		t.Fatalf("Expected to succeed, but failed: %q with message: %q", testcli.Error(), testcli.Stderr())
	}
	// Assert output
	for _, output := range []string{
		"Removing all modules in vendor/modules",
		"[terrafile-test-tag] Fetching git::https://github.com/terraform-aws-modules/terraform-aws-eks.git?ref=v18.29.0",
		"[terrafile-test-https] Fetching git::https://github.com/terraform-aws-modules/terraform-aws-eks.git",
		"[terrafile-test-registry] Found module version 18.29.0 at registry.terraform.io",
		"[terrafile-test-registry] Downloading from source URL git::https://github.com/terraform-aws-modules/terraform-aws-eks?ref=v18.29.0",
		"[terrafile-test-registry] Fetching git::https://github.com/terraform-aws-modules/terraform-aws-eks?ref=v18.29.0",
	} {
		assert.Contains(t, testcli.Stdout(), output)
	}
	// Assert files exist
	for _, moduleName := range []string{
		"terrafile-test-registry",
		"terrafile-test-https",
		"terrafile-test-tag",
	} {
		assert.DirExists(t, path.Join(workingDirectory, "vendor/modules", moduleName))
	}
}

func setup(t *testing.T) (current string, back func()) {
	folder, err := ioutil.TempDir("", "")
	assert.NoError(t, err)
	createTerrafile(t, folder)
	return folder, func() {
		assert.NoError(t, os.RemoveAll(folder))
	}
}

func createFile(t *testing.T, filename string, contents string) {
	assert.NoError(t, ioutil.WriteFile(filename, []byte(contents), 0644))
}

func createTerrafile(t *testing.T, folder string) {
	var yaml = `terrafile-test-registry:
  source: "terraform-aws-modules/eks/aws"
  version: "18.29.x"
terrafile-test-https:
  source: "https://github.com/terraform-aws-modules/terraform-aws-eks.git"
terrafile-test-tag:
  source: "https://github.com/terraform-aws-modules/terraform-aws-eks.git"
  version: "v18.29.0"
`
	createFile(t, path.Join(folder, "Terrafile.test"), yaml)
}
