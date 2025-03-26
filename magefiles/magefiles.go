// © 2019-2024 Diarkis Inc. All rights reserved.

//go:build mage
// +build mage

package main

import (
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

type Examples mg.Namespace

// Init Initialize a new project from the template, parameters: project_id builder_token output module_name
//
// Example of use
// init 012345678 sampleToken /tmp/sample-project github.com/sample-organization/sample-project
func Init(projectID, builderToken, output, moduleName string) error {
	return sh.RunV("go", "run", "./tools/init.go", projectID, builderToken, output, moduleName)
}

// Install Copy the examples to a new directory, parameters: project_id builder_token output
//
// Example of use
// install 012345678 sampleToken /tmp/examples-project
func (Examples) Install(projectID, builderToken, output string) error {
	return sh.RunV("go", "run", "./tools/install_examples.go", projectID, builderToken, output)
}
