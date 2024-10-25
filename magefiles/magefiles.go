//go:build mage
// +build mage

package main

import (
	"github.com/magefile/mage/sh"
)

// Init Initialize a new project from the template, parameters: project_id builder_token output module_name
//
// Example of use
// init 012345678 sampleToken /tmp/sample-project github.com/sample-organization/sample-project
func Init(projectID, builderToken, output, moduleName string) error {
	return sh.RunV("go", "run", "./tools/init.go", projectID, builderToken, output, moduleName)
}
