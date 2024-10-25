// © 2019-2024 Diarkis Inc. All rights reserved.

//go:build ignore

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	os.Exit(run())
}

func run() int {
	if len(os.Args) != 4 && len(os.Args) != 5 {
		usage()
	}

	projectID := os.Args[1]
	builderToken := os.Args[2]
	outputPath := os.Args[3]

	var moduleName string
	if len(os.Args) == 5 {
		moduleName = os.Args[4]
	}

	if moduleName == "" {
		moduleName = filepath.Base(outputPath)
	}

	fmt.Printf("- projectID: %q\n", projectID)
	fmt.Printf("- builderToken: %q\n", builderToken)
	fmt.Printf("- outputPath: %q\n", outputPath)
	fmt.Printf("- moduleName: %q\n", moduleName)

	err := checkGoModuleName(moduleName)
	if err != nil {
		fmt.Printf("check go module error: %v\n", err)
		return 1
	}

	err = runInstall(projectID, builderToken, outputPath)
	if err != nil {
		fmt.Printf("install error: %v\n", err)
		return 1
	}

	err = runRewriteImport(outputPath, moduleName)
	if err != nil {
		fmt.Printf("rewrite import error: %v\n", err)
		return 1
	}

	err = runPostSetup(outputPath, moduleName)
	if err != nil {
		fmt.Printf("post setup error: %v\n", err)
		return 1
	}

	return 0
}

func usage() {
	basename := filepath.Base(os.Args[0])
	fmt.Printf(`%[1]s is a tools to generate diarkis project.

Usage:
        %[1]s projectID builderToken outputPath <moduleName optional>
Sample:
        %[1]s 012345678 sampleToken /tmp/sample-project
        or
        %[1]s 012345678 sampleToken /tmp/sample-project github.com/sample-organization/sample-project`,
		basename)
	fmt.Printf("\n\n")

	os.Exit(1)
}

func checkGoModuleName(name string) error {
	// check module name using go mod init
	tempDir, err := os.MkdirTemp("", "diarkis-server-template")
	if err != nil {
		fmt.Printf("fail to create temporary directory. %v\n", err)
		return err
	}
	defer os.RemoveAll(tempDir)

	cmd := commandV("go", "mod", "init", name)
	cmd.Dir = tempDir

	return cmd.Run()
}

func runInstall(projectID, builderToken, outputPath string) error {
	args := []string{
		"run",
		"./tools/install",
		projectID,
		builderToken,
		outputPath,
	}
	cmd := commandV("go", args...)

	return cmd.Run()
}

func runRewriteImport(outputPath, moduleName string) error {
	args := []string{
		"run",
		"./tools/rewrite_import.go",
		outputPath,
		"github.com/Diarkis/diarkis-server-template",
		moduleName,
	}
	cmd := commandV("go", args...)

	return cmd.Run()
}

func runPostSetup(outputPath, moduleName string) error {
	args := []string{"mod", "edit", "-module", moduleName}
	cmd := commandV("go", args...)
	cmd.Dir = outputPath

	if err := cmd.Run(); err != nil {
		return err
	}

	// TODO(Henry): once merged remove the use of "vars.sh" and replace it
	// with a call to go mod edit -json | jq -r '.Module.Path'
	varshFilename := filepath.Join(outputPath, "puffer", "vars.sh")
	data := fmt.Sprintf("PROJECT_NAME=%s\n", moduleName)
	return os.WriteFile(varshFilename, []byte(data), 0666)
}

func commandV(name string, arg ...string) *exec.Cmd {
	cmd := exec.Command(name, arg...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd
}
