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

// List list available example
func (Examples) List() error {
	return sh.RunV("go", "run", "./tools/add_example.go", "-list")
}

// Add add a specific sample to an existing project
func (Examples) Add(name, projectDir string) error {
	return sh.RunV("go", "run", "./tools/add_example.go", "-add", "-destDir", projectDir, name)
}
