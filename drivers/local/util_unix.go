//go:build !windows

package local

import (
	"errors"
	"io/fs"
	"strings"

	"golang.org/x/sys/unix"
)

func isHidden(f fs.FileInfo, _ string) bool {
	return strings.HasPrefix(f.Name(), ".")
}

func isCrossDeviceError(err error) bool {
	return errors.Is(err, unix.EXDEV)
}
