package creation

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/abenz1267/related/config"
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

func getTemplate(name string, fragment config.Fragment) (bytes.Buffer, error) {
	var buffer bytes.Buffer

	tmplFile := strings.ReplaceAll(fragment.Template, "/", string(filepath.Separator))

	path := filepath.Join(".related", "templates", tmplFile)

	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		cfgDir, err := os.UserConfigDir()
		if err != nil {
			return buffer, fmt.Errorf("couldn't execute template: %w", err)
		}

		path = filepath.Join(cfgDir, "related", "templates", tmplFile)

		_, err = os.Stat(path)
		if err != nil {
			return buffer, fmt.Errorf("couldn't execute template: %w", err)
		}
	}

	wDir, err := os.Getwd()
	if err != nil {
		return buffer, nil
	}

	data := map[string]string{
		"wdir": wDir,
		"path": filepath.Dir(fragment.File),
		"name": name,
		"ext":  strings.Join(strings.Split(filepath.Base(fragment.File), ".")[1:], ""),
	}

	tmpl, err := template.ParseFiles(path)
	if err != nil {
		return buffer, fmt.Errorf("couldn't execute template: %w", err)
	}

	err = tmpl.Execute(&buffer, &data)
	if err != nil {
		return buffer, fmt.Errorf("couldn't execute template: %w", err)
	}

	return buffer, nil
}
