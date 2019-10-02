// Copyright © 2019 Tim Birkett <tim.birkett@devopsmakers.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package xterrafile

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/blang/semver"
	"github.com/hashicorp/terraform/registry"
	"github.com/hashicorp/terraform/registry/regsrc"
	"github.com/hashicorp/terraform/svchost/disco"

	jww "github.com/spf13/jwalterweatherman"
)

// IsRegistrySourceAddr check an address is a valid registry address
func IsRegistrySourceAddr(addr string) bool {
	jww.DEBUG.Printf("Testing if %s is a registry source", addr)
	_, err := regsrc.ParseModuleSource(addr)
	return err == nil
}

// GetRegistrySource retrieves a modules download source from a Terraform registry
func GetRegistrySource(name string, source string, version string) bool {
	modSrc, err := getModSrc(source)
	CheckIfError(name, err)

	version, err = getRegistryVersion(modSrc, version, nil)
	CheckIfError(name, err)

	jww.ERROR.Printf("[%s] Found module version %s at %s", name, version, modSrc.Host())

	return true
}

// Helper function to return a valid version
func getRegistryVersion(modSrc *regsrc.Module, version string, services *disco.Disco) (string, error) {
	// Don't log from Terraform's HTTP client
	log.SetFlags(0)
	log.SetOutput(ioutil.Discard)

	regClient := registry.NewClient(services, nil)
	regClientResp, err := regClient.ModuleVersions(modSrc)
	if err != nil {
		return "", err
	}

	validModuleVersionRange, err := semver.ParseRange(version)
	if err != nil {
		return "", err
	}

	regModule := regClientResp.Modules[0]
	for _, moduleVersion := range regModule.Versions {
		v, _ := semver.ParseTolerant(moduleVersion.Version)

		if validModuleVersionRange(v) {
			return v.String(), nil
		}
	}
	err = fmt.Errorf(
		"Unable to find a valid version at %s newest version is %s",
		modSrc.Host(),
		regModule.Versions[0].Version)
	return "", err
}

// Helper function to parse and return a module source
func getModSrc(source string) (*regsrc.Module, error) {
	modSrc, err := regsrc.ParseModuleSource(source)
	if err != nil {
		return nil, err
	}
	return modSrc, nil
}
