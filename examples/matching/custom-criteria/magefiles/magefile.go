// © 2019-2024 Diarkis Inc. All rights reserved.

//go:build mage
// +build mage

package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
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

const rootDir = "../.."

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
type Diarkis mg.Namespace

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

// Windows Build server binary for windows use
func (Build) Windows() error {
	fmt.Println("Build server binaries")

	return build("./build/windows-build.yml")
}

// Server Start a server locally.
//
// target can be either mars or http
func Server(target string) error {
	var exe string
	var args []string

	if target == "mars" {
		exe = filepath.Join("remote_bin", "mars"+binExt)
		args = append(args, filepath.Join("configs", "mars", "main.json"))
	} else {
		exe = filepath.Join(currDir, "remote_bin", target+binExt)
		// do not allow http to use all the cpu by default
		args = []string{"-c", "1"}
	}

	fmt.Printf("Starting %s server...\n", target)

	return sh.RunV(exe, args...)
}

// Coderefs Download the coderefs associated to the diarkis version in the go.mod.
func (Diarkis) Coderefs() error {
	version, err := getDiarkisVersion()
	if err != nil {
		fmt.Printf("fail to retrieve diarkis version from go.mod. %v", err)
		return err
	}

	return diarkisChangeVersion(version)
}

func diarkisChangeVersion(version string) error {
	fmt.Printf("version: %s\n", version)

	// check if the version exists
	url := fmt.Sprintf("https://docs.diarkis.io/sdk/server_coderef/%s.tar.gz", version)
	resp, err := http.Head(url)
	if err != nil {
		return fmt.Errorf("error cannot check the coderefs version: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("cannot find diarkis version %s\n", version)
		return errors.New("version not found")
	}

	resp, err = http.Get(url)
	if err != nil {
		return fmt.Errorf("error cannot download the coderefs: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("cannot find diarkis version %s\n", version)
		return errors.New("version not found")
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error cannot read request body: %w", err)
	}

	err = sh.Rm(filepath.Join("coderefs", version))
	if err != nil {
		return err
	}
	err = extractDocRefs(data, "coderefs")
	if err != nil {
		return fmt.Errorf("error cannot extract coderefs: %w", err)
	}

	err = sh.RunV("go", "mod", "edit",
		"-replace", "github.com/Diarkis/diarkis=./coderefs/"+version,
	)

	return err
}

func extractDocRefs(data []byte, dest string) error {
	gzipReader, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer gzipReader.Close()

	tarReader := tar.NewReader(gzipReader)
	for {
		tarHeader, err := tarReader.Next()
		if err == io.EOF {
			break
		}

		if tarHeader.Typeflag != tar.TypeDir && tarHeader.Typeflag != tar.TypeReg {
			fmt.Printf("skip non regular file %s (type %v)\n", tarHeader.Name, tarHeader.Typeflag)
			continue
		}
		// fmt.Printf("tarHeader.Name: %s\n", tarHeader.Name)
		abspath := filepath.Join(dest, tarHeader.Name)
		if tarHeader.Typeflag == tar.TypeDir {
			// directory
			fmt.Printf("create directory %s\n", tarHeader.Name)
			err = os.MkdirAll(abspath, tarHeader.FileInfo().Mode())
			if err != nil {
				return err
			}
			continue
		}

		fmt.Printf("extract file %s\n", tarHeader.Name)
		f, err := os.OpenFile(abspath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, tarHeader.FileInfo().Mode())
		if err != nil {
			return err
		}
		_, err = io.Copy(f, tarReader)
		f.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

func getDiarkisVersion() (string, error) {
	stdout := bytes.NewBuffer(nil)
	cmd := exec.Command("go", "mod", "edit", "-json")
	cmd.Stdout = stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return "", err
	}

	var doc struct {
		Module struct {
			Path string
		}
		Go      string
		Require []struct {
			Path    string
			Version string
		}
	}
	err = json.Unmarshal(stdout.Bytes(), &doc)
	if err != nil {
		return "", err
	}

	for _, require := range doc.Require {
		if require.Path == "github.com/Diarkis/diarkis" {
			return require.Version, nil
		}
	}

	return "", errors.New("diarkis is not in the requirement of the go mod file")
}

// build invoke diarkis-cli with the proper configuration file.
func build(buildCfg string) error {
	err := cleanDir("remote_bin")
	if err != nil {
		return err
	}

	diarkisCli := filepath.Join(getProjectRoot(), getDiarkisCli())

	return runVInDir(currDir, diarkisCli, "build", "-c", buildCfg, "--host", diarkisCLIHost)
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
	return filepath.Clean(rootDir)
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
