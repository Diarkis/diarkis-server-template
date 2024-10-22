// © 2019-2024 Diarkis Inc. All rights reserved.

//go:build unix

package main

import (
	"fmt"
	"io/fs"
	"os"
	"syscall"
)

func applyOwnership(sourcePath, destPath string, fileInfo fs.FileInfo) error {
	stat, ok := fileInfo.Sys().(*syscall.Stat_t)
	if !ok {
		return fmt.Errorf("Failed to get raw syscall.Stat_t data for \x1b[0;91m %v \x1b[0m", sourcePath)
	}

	return os.Lchown(destPath, int(stat.Uid), int(stat.Gid))
}
