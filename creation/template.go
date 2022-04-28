package creation

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	PermFile   = 0o644
	PermFolder = 0o755
)

func Template(result bool, global bool, src, dest string) error {
	destFolder := filepath.Join(".related", "templates")

	if global {
		userConfig, err := os.UserConfigDir()
		if err != nil {
			return fmt.Errorf("couldn't get user config dir: %w", err)
		}

		destFolder = filepath.Join(userConfig, "related", "templates")
	}

	source, err := os.ReadFile(src)
	if err != nil {
		return fmt.Errorf("couldn't read source file: %w", err)
	}

	dest = filepath.Join(destFolder, dest+".tmpl")

	err = os.MkdirAll(filepath.Dir(dest), PermFolder)
	if err != nil {
		return fmt.Errorf("couldn't create destination folder(s): %w", err)
	}

	err = os.WriteFile(dest, source, PermFile)
	if err != nil {
		return fmt.Errorf("couldn't write template file: %w", err)
	}

	if result {
		fmt.Println(dest) //nolint
	}

	return nil
}
