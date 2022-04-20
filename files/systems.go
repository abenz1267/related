package files

import (
	"errors"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/abenz1267/gonerics"
)

const (
	ConfigDir   SystemDir = "related"
	ProjectDir  SystemDir = ".related"
	TemplateDir TypeDir   = "templates"
	ScriptDir   TypeDir   = "scripts"
)

type (
	TypeDir   string
	SystemDir string
)

var configDir SystemDir //nolint

func init() { //nolint
	configDir = userConfigDir()
}

func SetConfigDir(dir string) {
	configDir = SystemDir(dir)
}

func configFS() fs.FS {
	path := filepath.Join(string(configDir), string(ConfigDir))

	if exists(path) {
		return os.DirFS(path)
	}

	return nil
}

func projectFS() fs.FS {
	if exists(string(ProjectDir)) {
		return os.DirFS(string(ProjectDir))
	}

	return nil
}

func exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false
		}

		log.Panic(err)
	}

	return true
}

func Systems() []fs.FS {
	return []fs.FS{projectFS(), configFS()}
}

func userConfigDir() SystemDir {
	dir := gonerics.TryResult(os.UserConfigDir())

	return SystemDir(dir)
}
