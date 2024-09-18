// © 2019-2024 Diarkis Inc. All rights reserved.

package main

import (
	"fmt"
	"go/format"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) != 4 {
		fmt.Printf("usage %s: <target_dir> <oldImport> <newImport>\n", os.Args[0])
		os.Exit(1)
	}

	targetDir := os.Args[1]
	oldImport := os.Args[2]
	newImport := os.Args[3]
	err := processDir(targetDir, oldImport, newImport)
	if err != nil {
		os.Exit(1)
	}
}

func processDir(dir, oldImport, newImport string) error {
	err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		if info.IsDir() {
			return nil
		}

		ext := filepath.Ext(info.Name())
		if ext != ".go" {
			return nil
		}

		fixedStr, rewrote, err := fixImport(path, oldImport, newImport)
		if err != nil {
			return err
		}
		if rewrote {
			err = os.WriteFile(path, []byte(fixedStr), info.Mode())
			if err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

func fixImport(filename, oldImport, newImport string) (string, bool, error) {
	fset := token.NewFileSet()

	expr, err := parser.ParseFile(fset, filename, nil, parser.Mode(0))
	if err != nil {
		return "", false, err
	}

	found := false
	for _, imp := range expr.Imports {
		if strings.HasPrefix(imp.Path.Value[1:], oldImport) {
			// rewrite import
			imp.Path.Value = strings.Replace(imp.Path.Value, oldImport, newImport, 1)
			found = true
		}
	}

	var buf strings.Builder

	// write into buf
	if err := format.Node(&buf, token.NewFileSet(), expr); err != nil {
		return "", false, err
	}

	return buf.String(), found, nil
}
