// © 2019-2024 Diarkis Inc. All rights reserved.

package main

import (
	"fmt"
	"os"
	"unicode/utf8"

	"io"
	"path/filepath"
	"strings"
	"syscall"
)

var projectID = ""
var buildToken = ""

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Path failed \x1b[0;91m %v \x1b[0m\n", err)
		os.Exit(1)
		return
	}
	src := fmt.Sprintf("%s/src/", cwd)
	dest := ""
	projectID = os.Args[1]
	buildToken = os.Args[2]
	if os.Args[3][0:1] == "/" {
		dest = os.Args[3]
	} else {
		dest = filepath.Join(cwd, os.Args[3])
	}
	list := strings.Split(dest, "/")
	pkg := list[len(list)-1]
	fmt.Printf("\x1b[0;90m Installing the template as package %s to %s \x1b[0m\n", pkg, dest)
	_, err = os.Stat(dest)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.Mkdir(dest, os.FileMode(0777))
			if err != nil {
				fmt.Printf("Error \x1b[0;91m %v \x1b[0m\n", err)
				os.Exit(1)
			}
		} else {
			fmt.Printf("Error \x1b[0;91m %v \x1b[0m\n", err)
			os.Exit(1)
		}
	}
	err = copyDirectory(pkg, src, dest)
	if err != nil {
		fmt.Printf("Error \x1b[0;91m %v \x1b[0m\n", err)
		os.Exit(1)
	}
	fmt.Printf("Installation of template completed - \x1b[0;32m %v \x1b[0m\n", dest)
	err = os.Chdir(dest)
	if err != nil {
		fmt.Printf("Error \x1b[0;91m %v \x1b[0m\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}

func copyDirectory(pkg string, src string, dest string) error {
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		sourcePath := filepath.Join(src, entry.Name())
		destPath := filepath.Join(dest, entry.Name())
		fileInfo, err := os.Stat(sourcePath)
		if err != nil {
			return err
		}
		if entry.Name() == "build.yml" {
			_, err := os.Stat(destPath)
			if err == nil {
				fmt.Printf("Project build.yml found at %s - Skip installing\n", destPath)
				// the project already has build.yml - skip
				continue
			}
		}
		if entry.Name() == "local-build.yml" {
			_, err := os.Stat(destPath)
			if err == nil {
				fmt.Printf("Project local-build.yml found at %s - Skip installing\n", destPath)
				// the project already has local-build.yml - skip
				continue
			}
		}
		if entry.Name() == "go.mod" {
			_, err := os.Stat(destPath)
			if err == nil {
				fmt.Printf("Project go.mod found at %s - Skip installing\n", destPath)
				// the project already has build.yml - skip
				continue
			}
		}
		stat, ok := fileInfo.Sys().(*syscall.Stat_t)
		if !ok {
			return fmt.Errorf("Failed to get raw syscall.Stat_t data for \x1b[0;91m %v \x1b[0m", sourcePath)
		}
		switch fileInfo.Mode() & os.ModeType {
		case os.ModeDir:
			if err := createIfNotExists(destPath, 0755); err != nil {
				return err
			}
			if err := copyDirectory(pkg, sourcePath, destPath); err != nil {
				return err
			}
		case os.ModeSymlink:
			if err := copySymLink(sourcePath, destPath); err != nil {
				return err
			}
		default:
			if err := copyFile(pkg, sourcePath, destPath); err != nil {
				return err
			}
		}

		if err := os.Lchown(destPath, int(stat.Uid), int(stat.Gid)); err != nil {
			return err
		}
		info, err := entry.Info()
		if err != nil {
			fmt.Printf("\x1b[0;91mFailed to get file info from %s\x1b[0m\n", entry.Name())
			continue
		}

		isSymlink := info.Mode()&os.ModeSymlink != 0
		if !isSymlink {
			if err := os.Chmod(destPath, info.Mode()); err != nil {
				return err
			}
		}
	}
	return nil
}

func copyFile(pkg string, srcFile string, dstFile string) error {
	fmt.Printf("\x1b[0;90m Installing from %v to %v \x1b[0m\n", srcFile, dstFile)
	out, err := os.Create(dstFile)
	if err != nil {
		return err
	}
	defer out.Close()
	in, err := os.Open(srcFile)
	defer in.Close()
	if err != nil {
		return err
	}
	data, err := io.ReadAll(in)
	if err != nil {
		return err
	}
	fileData := string(data)
	if utf8.ValidString(string(data)) {
		fileData = strings.Replace(fileData, "{{PROJECT_ID}}", projectID, -1)
		fileData = strings.Replace(fileData, "{{BUILD_TOKEN}}", buildToken, -1)
	} else {
		fmt.Printf("\x1b[38;5;220mBinary file detected, skipping the replace. %s\x1b[0m\n", srcFile)
	}
	_, err = io.WriteString(out, fileData)
	if err != nil {
		return err
	}
	in.Sync()
	os.Chmod(dstFile, 0700)
	return nil
}

func exists(filePath string) bool {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	}
	return true
}

func createIfNotExists(dir string, perm os.FileMode) error {
	if exists(dir) {
		return nil
	}
	if err := os.MkdirAll(dir, perm); err != nil {
		return fmt.Errorf("Failed to create directory: '%s', Error:\x1b[0;91m %v \x1b[0m", dir, err.Error())
	}

	return nil
}

func copySymLink(source, dest string) error {
	link, err := os.Readlink(source)
	if err != nil {
		return err
	}
	return os.Symlink(link, dest)
}
