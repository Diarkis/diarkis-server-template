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
	"strings"

	// mg contains helpful utility functions, like Deps
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

type Diarkis mg.Namespace
type Build mg.Namespace
type Puffer mg.Namespace

// Version Print the version of diarkis currently used.
func (Diarkis) Version() error {
	version, err := getDiarkisVersion()
	if err != nil {
		return err
	}
	fmt.Printf("diarkis version is %s\n", version)
	return nil
}

// ChangeVersion Change diarkis version.
func (Diarkis) ChangeVersion(version string) error {
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

	// update go mod
	err = updateGoMod(version)
	if err != nil {
		return err
	}
	return nil
}

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

// Server Start a server locally: [ target=mars ] [ target=http ] [ target=udp ] [ target=tcp ]
func Server(target string) error {
	var exe string
	var args []string
	var exeExt string

	if runtime.GOOS == "windows" {
		exeExt = ".exe"
	}

	if target == "mars" {
		exe = filepath.Join("remote_bin", "mars"+exeExt)
		args = append(args, "./configs/mars/main.json")
	} else {
		exe = filepath.Join("remote_bin", target+exeExt)
	}

	fmt.Printf("Starting %s server...\n", target)

	return sh.RunV(exe, args...)
}

// Gen Generate go, cpp, and cs code files using puffer (Diarkis packet gen module) from packet definition written in json.
func (Puffer) Gen() error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("fail to retrieve current working directory: %w", err)
	}
	// read project name
	data, err := os.ReadFile(filepath.Join("puffer", "vars.sh"))
	if err != nil {
		return err
	}
	lines := strings.SplitN(string(data), "\n", 1)
	projectName, ok := strings.CutPrefix(lines[0], "PROJECT_NAME=")

	if !ok {
		err = errors.New("puffer/vars.sh not well formed")
		fmt.Printf("%v\n", err)
		return err
	}

	var pufferBin string
	switch runtime.GOOS {
	case "linux":
		pufferBin = filepath.Join("puffer", "puffer-linux")
	case "darwin":
		pufferBin = filepath.Join("puffer", "puffer-mac")
	case "windows":
		// TODO
	}

	if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
		pufferBin = filepath.Join(cwd, pufferBin)
	}

	cmd := exec.Command(pufferBin, ".", ".", projectName)
	cmd.Dir = "puffer"
	cmd.Stdout = os.Stdout

	return cmd.Run()
}

// Clean Delete all generated protocol code files.
func (Puffer) Clean() error {
	for _, dir := range []string{"go", "cs", "cpp"} {
		err := sh.Rm(filepath.Join("puffer", dir))
		if err != nil {
			return err
		}
	}
	return nil
}

// GoCli Starts Go test client host=<HTTP address> uid=<client user ID> clientKey=<client key> puffer=<true/false>
func GoCli(host, uid, clientKey, puffer string) error {
	bin := filepath.Join("remote_bin", "testcli")
	if runtime.GOOS == "windows" {
		bin += ".exe"
	}
	return sh.RunV(bin,
		fmt.Sprintf("--host=%s", host),
		fmt.Sprintf("--uid=%s", uid),
		fmt.Sprintf("--clientKey=%s", clientKey),
		fmt.Sprintf("--puffer=%s", puffer),
	)
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

func build(buildCfg string) error {
	var diarkisCli string
	switch runtime.GOOS {
	case "linux":
		diarkisCli = "./diarkis-cli/os/linux/bin/diarkis-cli"
	case "darwin":
		diarkisCli = "./diarkis-cli/os/mac/bin/diarkis-cli"
	case "windows":
		diarkisCli = "./diarkis-cli/os/windows/bin/diarkis-cli.exe"
	}

	err := cleanDir("remote_bin")
	if err != nil {
		return err
	}

	// $(DIARKIS_CLI) build -c $(BUILD_CONFIG) --host v3.builder.diarkis.io
	return sh.RunV(diarkisCli, "build", "-c", buildCfg, "--host", "v3.builder.diarkis.io")
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

func updateGoMod(version string) error {
	err := sh.RunV("go", "mod", "edit",
		"-require", "github.com/Diarkis/diarkis@"+version,
		"-replace", "github.com/Diarkis/diarkis=./coderefs/"+version,
	)

	return err
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
