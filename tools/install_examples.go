// © 2019-2024 Diarkis Inc. All rights reserved.

//go:build ignore

package main

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/goccy/go-yaml"
)

func main() {
	if len(os.Args) != 4 {
		usage()
	}

	projectID := os.Args[1]
	buildToken := os.Args[2]
	dest := os.Args[3]

	err := copyExampleToTargetTemplate(projectID, buildToken, dest, ".")
	if err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}
	err = copyDiarkisCLIToTargetExample(dest, ".")
	if err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Installation done\n")
}

func usage() {
	basename := filepath.Base(os.Args[0])
	fmt.Printf(`%[1]s is a tools to install diarkis example project.

Usage:
        %[1]s projectID builderToken outputPath
Sample:
        %[1]s 012345678 sampleToken /tmp/sample-project`,
		basename)
	fmt.Printf("\n\n")

	os.Exit(1)
}

func copyExampleToTargetTemplate(projectID, builderToken, destDir, templateDir string) error {
	var yamlFiles []string

	fmt.Printf("Install examples to %s\n", destDir)

	examplesDir := filepath.Join(templateDir, "examples")
	magefilesDirectories := []string{}

	err := copyDir(examplesDir, destDir, func(relPath string) {
		// Keep track of each magefiles folder
		parentDir := filepath.Dir(relPath)
		if filepath.Base(relPath) == "magefile.go" && filepath.Base(parentDir) == "magefiles" {
			magefilesDirectories = append(magefilesDirectories, filepath.Join(destDir, parentDir))
		}
		if strings.EqualFold(filepath.Ext(relPath), ".yml") {
			yamlFiles = append(yamlFiles, filepath.Join(destDir, relPath))
		}
	})

	if err != nil {
		return err
	}

	// post process yaml file
	for _, abspath := range yamlFiles {
		err = updateBuildSettings(abspath, projectID, builderToken)
		if err != nil {
			return err
		}
	}

	// copy magefiles/common.go to each sample.
	magefilesDir := filepath.Join(templateDir, "examples", "magefiles")
	files, err := os.ReadDir(magefilesDir)
	if err != nil {
		return err
	}
	for _, file := range files {
		if file.Type().IsRegular() {
			stats, err := file.Info()
			if err != nil {
				return err
			}
			for _, dirname := range magefilesDirectories {
				destPath := filepath.Join(dirname, file.Name())

				err := copyFile(filepath.Join(magefilesDir, file.Name()), destPath, stats.Mode().Perm())
				if err != nil {
					return err
				}
			}
		}
	}

	// download code refs
	for _, dirname := range magefilesDirectories {
		parentDir := filepath.Dir(dirname)
		fmt.Printf("Download coderefs for %s\n", dirname)
		if err := runDiarkisCoderefs(parentDir); err != nil {
			return err
		}
	}

	return nil
}

func runDiarkisCoderefs(dirname string) error {
	cmd := exec.Command("go", "run", "./magefiles/mage.go", "diarkis:coderefs")
	cmd.Dir = dirname
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func copyDiarkisCLIToTargetExample(destDir, templateDir string) error {
	destDir = filepath.Join(destDir, "diarkis-cli")
	fmt.Printf("Copy diarkis-cli to %s\n", destDir)

	examplesDir := filepath.Join(templateDir, "src", "diarkis-cli")
	return copyDir(examplesDir, destDir, func(relPath string) {})
}

func copyDir(srcDir, destDir string, onCopy func(relPath string)) error {
	return filepath.Walk(srcDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}

		// compute relative path from srcDir
		rel, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}

		destPath := filepath.Join(destDir, rel)
		if info.IsDir() {
			err = os.MkdirAll(destPath, info.Mode().Perm())
			return err
		}

		if !info.Mode().IsRegular() {
			fmt.Printf("skip non regular file %s (%s)\n", rel, info.Mode())
			return nil
		}

		onCopy(rel)

		return copyFile(path, destPath, info.Mode().Perm())
	})
}

func copyFile(srcPath, destPath string, perm fs.FileMode) error {
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	dstFile, err := os.OpenFile(destPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, perm)
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
