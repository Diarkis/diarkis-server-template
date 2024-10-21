// © 2019-2024 Diarkis Inc. All rights reserved.

//go:build ignore

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const copyRightHeader = "Diarkis Inc. All rights reserved."

func main() {
	os.Exit(run())
}

func run() int {
	year := time.Now().Year()
	headerWithYear := fmt.Sprintf("// © 2019-%d %s\n\n", year, copyRightHeader)

	err := filepath.Walk(".", func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}

		if !info.Mode().IsRegular() {
			return nil
		}

		ext := filepath.Ext(info.Name())
		if !strings.EqualFold(ext, ".go") {
			return nil
		}

		found, err := hasCopyrightHeader(path)
		if err != nil {
			return err
		}
		if found {
			return nil
		}

		err = addCopyrightHeader(path, headerWithYear)
		return err
	})

	if err != nil {
		return 1
	}

	return 0
}

func hasCopyrightHeader(filename string) (bool, error) {
	f, err := os.Open(filename)
	if err != nil {
		return false, err
	}
	defer f.Close()

	fileScanner := bufio.NewScanner(f)

	fileScanner.Split(bufio.ScanLines)

	// check the first line only
	for fileScanner.Scan() {
		if strings.Contains(fileScanner.Text(), copyRightHeader) {
			return true, nil
		}
		return false, nil
	}

	return false, nil
}

func addCopyrightHeader(filename string, headerWithYear string) error {
	// Try to keep the same permissions.
	stat, err := os.Stat(filename)
	if err != nil {
		return err
	}
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	b := bytes.NewBuffer(nil)
	b.Grow(len(data) + len(headerWithYear))
	_, _ = b.WriteString(headerWithYear)
	_, _ = b.Write(data)

	// write to a temp file and rename it.
	// same as what did the Makefile previously.
	const tempFilename = "temp"
	err = os.WriteFile(tempFilename, b.Bytes(), stat.Mode())
	if err != nil {
		return err
	}

	return os.Rename(tempFilename, filename)
}
