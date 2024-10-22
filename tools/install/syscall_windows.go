// © 2019-2024 Diarkis Inc. All rights reserved.

//go:build windows

package main

import "io/fs"

func applyOwnership(sourcePath, destPath string, fileInfo fs.FileInfo) error {
	return nil
}
