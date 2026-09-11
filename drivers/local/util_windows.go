//go:build windows

package local

import (
	"errors"
	"io/fs"
	"path/filepath"
	"syscall"

	"github.com/OpenListTeam/OpenList/v4/internal/model"
	"golang.org/x/sys/windows"
)

func isHidden(f fs.FileInfo, fullPath string) bool {
	filePath := filepath.Join(fullPath, f.Name())
	namePtr, err := syscall.UTF16PtrFromString(filePath)
	if err != nil {
		return false
	}
	attrs, err := syscall.GetFileAttributes(namePtr)
	if err != nil {
		return false
	}
	return attrs&syscall.FILE_ATTRIBUTE_HIDDEN != 0
}

func isCrossDeviceError(err error) bool {
	return errors.Is(err, windows.ERROR_NOT_SAME_DEVICE)
}
