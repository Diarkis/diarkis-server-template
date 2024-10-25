// © 2019-2024 Diarkis Inc. All rights reserved.

//go:build mage
// +build mage

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	// mg contains helpful utility functions, like Deps
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

var binExt string
var currDir string

const diarkisCLIHost = "v3.builder.diarkis.io"

func init() {
	if runtime.GOOS == "windows" {
		binExt = ".exe"
	}

	cd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	currDir = cd
}

// Default target to run when none is specified
// If not set, running mage will list available targets
// var Default = Build

type Build mg.Namespace

// Local Build server binary for local use
func (Build) Local() error {
	fmt.Println("Build server binaries for local use")

	var buildCfg string
	switch runtime.GOOS {
	case "linux":
		buildCfg = "./build/linux-build.yml"
	case "darwin":
		buildCfg = "./build/mac-build.yml"
	case "windows":
		buildCfg = "./build/windows-build.yml"
	}

	return build(buildCfg)
}

// Linux Build server binary for linux or container environment
func (Build) Linux() error {
	fmt.Println("Build server binaries")

	return build("./build/linux-build.yml")
}

// Mac Build server binary for mac use
func (Build) Mac() error {
	fmt.Println("Build server binaries")

	return build("./build/mac-build.yml")
}

// Server Start a server locally: [ target=mars ] [ target=http ]
func Server(target string) error {
	var exe string
	var args []string

	projectRoot := getProjectRoot()

	if target == "mars" {
		exe = filepath.Join("remote_bin", "mars"+binExt)
		args = append(args, filepath.Join("configs", "mars", "main.json"))
	} else {
		exe = filepath.Join(currDir, "remote_bin", target+binExt)
	}

	fmt.Printf("Starting %s server...\n", target)

	return runVInDir(projectRoot, exe, args...)
}

func build(buildCfg string) error {
	err := cleanDir("remote_bin")
	if err != nil {
		return err
	}

	diarkisCli := filepath.Join(getProjectRoot(), getDiarkisCli())

	if err = buildDependencies(buildCfg); err != nil {
		return err
	}

	return runVInDir(currDir, diarkisCli, "build", "-c", buildCfg, "--host", diarkisCLIHost)
}

// buildDependencies build project root artifacts
func buildDependencies(buildCfg string) error {
	fmt.Printf("buildDependencies\n")
	projectRoot := getProjectRoot()
	diarkisCli := getDiarkisCli()

	// build the dependencies with the same configuration file name.
	// it should exist at the project root directory
	return runVInDir(projectRoot, diarkisCli, "build", "-c", buildCfg, "--host", diarkisCLIHost)
}

func getDiarkisCli() string {
	switch runtime.GOOS {
	case "linux":
		return "./diarkis-cli/os/linux/bin/diarkis-cli"
	case "darwin":
		return "./diarkis-cli/os/mac/bin/diarkis-cli"
	case "windows":
		return `diarkis-cli\os\windows\bin\diarkis-cli.exe`
	default:
		panic("unsupported platform")
	}
}

func getProjectRoot() string {
	return filepath.Clean("../../..")
}

func cleanDir(dir string) error {
	files, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	for _, file := range files {
		abspath := filepath.Join(dir, file.Name())
		if file.IsDir() {
			err = sh.Rm(abspath)
			if err != nil {
				return err
			}
			continue
		}
		err = os.Remove(abspath)
		if err != nil {
			return err
		}
	}

	return nil
}

// runVInDir mostly the same as sh.RunV but run the command inside dir
func runVInDir(dir, cmd string, args ...string) error {
	c := exec.Command(cmd, args...)
	c.Env = os.Environ()
	c.Stderr = os.Stderr
	c.Stdout = os.Stdout
	c.Stdin = os.Stdin
	c.Dir = dir

	return c.Run()
}
