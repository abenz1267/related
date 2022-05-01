package creation

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/abenz1267/related/config"
)

func Create(kind string, name string, file string, result bool) {
	cfg, err := config.Get()
	if err != nil {
		panic(err)
	}

	fragment, err := cfg.GetFragment(name, cfg.Fragments)
	if err != nil {
		panic(err)
	}

	fileName := strings.ReplaceAll(fragment.File, "<name>", file)

	data, err := getTemplate(file, fragment)
	if err != nil {
		panic(err)
	}

	if _, err := os.Stat(fileName); errors.Is(err, os.ErrNotExist) {
		err = os.MkdirAll(filepath.Dir(fileName), PermFolder)
		if err != nil {
			panic(err)
		}

		err = os.WriteFile(fileName, data.Bytes(), PermFile)
		if err != nil {
			panic(err)
		}
	}

	if result {
		fmt.Println(fileName) //nolint
	}
}
