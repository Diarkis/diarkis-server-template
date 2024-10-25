// © 2019-2024 Diarkis Inc. All rights reserved.

//go:build ignore

package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/goccy/go-yaml"
)

func main() {
	var listExample, addExample bool
	var destDir string

	flag.BoolVar(&listExample, "list", false, "List available example.")
	flag.BoolVar(&addExample, "add", false, "Add one example.")
	flag.StringVar(&destDir, "destDir", "", "Where to add the example")

	flag.Parse()

	fmt.Printf("destDir: %s\n", destDir)
	fmt.Printf("addExample: %t\n", addExample)
	fmt.Printf("listExample: %t\n", listExample)

	if listExample {
		l, err := getExampleList()
		if err != nil {
			fmt.Printf("fail to retrieve the example list. %v\n", err)
			os.Exit(1)
		}
		sort.Strings(l)
		for _, name := range l {
			fmt.Printf("- %s\n", name)
		}
		return
	}

	if addExample {
		for _, arg := range flag.Args() {
			ok, err := validateExample(arg)
			if err != nil {
				panic(err)
			}
			fmt.Printf("ok: %t\n", ok)

			err = copyExampleToTargetTemplate(destDir, "examples/"+arg, arg)
			if err != nil {
				panic(err)
			}
		}
		return
	}

	flag.Usage()
}

func validateExample(name string) (bool, error) {
	examples, err := getExampleList()
	if err != nil {
		return false, err
	}

	for _, ex := range examples {
		if ex == name {
			return true, nil
		}
	}

	return false, nil
}

func getExampleList() ([]string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	examplesDir := filepath.Join(cwd, "examples")
	stat, err := os.Stat(examplesDir)
	if err != nil {
		return nil, err
	}
	if !stat.IsDir() {
		return nil, errors.New("examples is expected to be a folder")
	}

	var examples []string
	err = filepath.Walk(examplesDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}

		if info.IsDir() {
			return nil
		}
		if info.Name() == "go.mod" {
			// found one example
			rel, err := filepath.Rel(examplesDir, filepath.Dir(path))
			if err != nil {
				return err
			}
			examples = append(examples, filepath.ToSlash(rel))
		}
		return nil
	})

	return examples, err
}

// getTargetTemplateDiarkisBuildInfo retrieve the diarkis.project_id
// and builder_token from an existing template.
func getTargetTemplateDiarkisBuildInfo(templateDir string) (string, string, error) {
	linuxBuildFile := filepath.Join(templateDir, "build", "linux-build.yml")
	fmt.Printf("linuxBuildFile: %q\n", linuxBuildFile)
	linuxData, err := os.ReadFile(linuxBuildFile)
	if err != nil {
		return "", "", err
	}

	var o struct {
		Diarkis struct {
			ProjectID    string `yaml:"project_id"`
			BuilderToken string `yaml:"builder_token"`
		} `yaml:"diarkis"`
	}
	err = yaml.Unmarshal(linuxData, &o)
	if err != nil {
		return "", "", err
	}

	return o.Diarkis.ProjectID, o.Diarkis.BuilderToken, nil
}

func copyExampleToTargetTemplate(templateDir, exampleDir, exampleName string) error {
	projectID, builderToken, err := getTargetTemplateDiarkisBuildInfo(templateDir)
	if err != nil {
		return err
	}

	exampleOutDir := filepath.Join(templateDir, "examples", exampleName)
	err = os.MkdirAll(exampleOutDir, 0755)
	if err != nil {
		return err
	}

	var yamlFiles []string
	err = filepath.Walk(exampleDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}

		// compute relative path from exampleDir
		rel, err := filepath.Rel(exampleDir, path)
		if err != nil {
			return err
		}

		if strings.EqualFold(filepath.Ext(info.Name()), ".yml") {
			yamlFiles = append(yamlFiles, rel)
		}

		destPath := filepath.Join(exampleOutDir, rel)
		if info.IsDir() {
			err = os.MkdirAll(destPath, info.Mode().Perm())
			return err
		}

		if !info.Mode().IsRegular() {
			return nil
		}

		err = func() error {
			srcFile, err := os.Open(path)
			if err != nil {
				return err
			}
			defer srcFile.Close()
			dstFile, err := os.OpenFile(destPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, info.Mode().Perm())
			if err != nil {
				return err
			}
			_, err = io.Copy(dstFile, srcFile)
			if err != nil {
				// best effort
				_ = dstFile.Close()
				_ = os.Remove(destPath)
				return err
			}

			return nil
		}()

		return err
	})

	if err != nil {
		return err
	}

	// post process yaml file
	for _, rel := range yamlFiles {
		destPath := filepath.Join(exampleOutDir, rel)
		err = updateBuildSettings(destPath, projectID, builderToken)
		if err != nil {
			return err
		}
	}

	// FIXME(Henry) it would be nice if we update the example go.mod
	// with the current version of diarkis used by the template project.
	// This would allow to have proper code completion

	return nil
}

func updateBuildSettings(filename, projectID, builderToken string) error {
	stat, err := os.Stat(filename)
	if err != nil {
		return err
	}
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	var m map[string]any
	err = yaml.Unmarshal(data, &m)
	if err != nil {
		// not the same format
		return nil
	}

	if _, ok := m["buildSettings"]; !ok {
		// not the same format
		return nil
	}

	diarkis, ok := m["diarkis"]
	if !ok {
		// not the same format
		return nil
	}

	diarkisMap, ok := diarkis.(map[string]any)
	if !ok {
		// not the same format
		return nil
	}

	for _, k := range []string{"project_id", "builder_token"} {
		v, ok := diarkisMap[k]
		if !ok {
			return nil
		}
		_, ok = v.(string)
		if !ok {
			return nil
		}
	}

	diarkisMap["project_id"] = projectID
	diarkisMap["builder_token"] = builderToken

	updated, err := yaml.Marshal(m)
	if err != nil {
		return err
	}

	return os.WriteFile(filename, updated, stat.Mode().Perm())
}
