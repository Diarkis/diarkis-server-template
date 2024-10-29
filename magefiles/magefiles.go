// © 2019-2024 Diarkis Inc. All rights reserved.

//go:build mage
// +build mage

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"

	"github.com/goccy/go-yaml"
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

func AnonymizeToken() error {
	err := filepath.Walk("examples", func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		if info.IsDir() {
			if filepath.Base(path) == "coderefs" {
				return filepath.SkipDir
			}
		} else {
			ext := filepath.Ext(path)
			if !strings.EqualFold(ext, ".yml") {
				return nil
			}

			return anonymizeBuildConfig(path)
		}
		return nil
	})

	return err
}

func anonymizeBuildConfig(filename string) error {
	stat, err := os.Stat(filename)
	if err != nil {
		return err
	}
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("fail to read file: %w", err)
	}

	var m map[string]any
	err = yaml.Unmarshal(data, &m)
	if err != nil {
		// not the same format
		return nil
	}

	d, ok := m["diarkis"]
	if !ok {
		return nil
	}
	dm, ok := d.(map[string]any)
	if !ok {
		return nil
	}
	_, ok = dm["builder_token"]
	if !ok {
		return nil
	}

	// Do not use yaml.Marshal because the map key order is not stable
	var sb strings.Builder

	scanner := bufio.NewScanner(bytes.NewBuffer(data))
	// optionally, resize scanner's capacity for lines over 64K, see next example
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(strings.TrimSpace(line), "builder_token") {
			tokens := strings.SplitN(line, ":", 2)
			line = fmt.Sprintf(`%s: "{{BUILD_TOKEN}}"`, tokens[0])
		} else if strings.HasPrefix(strings.TrimSpace(line), "project_id:") {
			tokens := strings.SplitN(line, ":", 2)
			line = fmt.Sprintf(`%s: "{{PROJECT_ID}}"`, tokens[0])
		}
		sb.WriteString(line)
		sb.WriteByte('\n')
	}

	return os.WriteFile(filename, []byte(sb.String()), stat.Mode().Perm())
}
