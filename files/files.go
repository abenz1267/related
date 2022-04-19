package files

import (
	"io/fs"
	"path/filepath"
	"runtime"
	"strings"
)

// FindFile finds a given template or script file within all available filesystems
// and returns the first result including the filesystem.
// Search order: project, user config, embedded.
func FindFile(name string, dir TypeDir) (string, fs.FS) {
	systems := Systems()

	path := NormalizeForFS(filepath.Join(string(dir), name))

	path = strings.ReplaceAll(path, "\\", "/")

	for _, system := range systems {
		if system == nil {
			continue
		}

		if _, err := system.Open(path); err == nil {
			return path, system
		}
	}

	return "", nil
}

func NormalizeForFS(in string) string {
	if runtime.GOOS == "windows" {
		return strings.ReplaceAll(in, "\\", "/")
	}

	return in
}
