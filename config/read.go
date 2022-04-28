package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/exp/slices"
	"gopkg.in/yaml.v2"
)

func Get() (Config, error) {
	var config Config

	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		return config, fmt.Errorf("error getting user config dir: %w", err)
	}

	folders := []string{filepath.Join(userConfigDir, "related"), ".related"}

	for _, v := range folders {
		if _, err = os.Stat(v); errors.Is(err, os.ErrNotExist) {
			continue
		}

		files, filesErr := getFiles(v)
		if filesErr != nil {
			return config, fmt.Errorf("error getting config files: %w", filesErr)
		}

		for _, f := range files {
			other, configErr := read(f)
			if configErr != nil {
				return config, fmt.Errorf("error reading config: %w", configErr)
			}

			config.merge(other)
		}
	}

	transformErr := config.transform()
	if transformErr != nil {
		return config, transformErr
	}

	config.sort()

	return config, nil
}

func read(path string) (Config, error) {
	var config Config

	b, err := os.ReadFile(path)
	if err != nil {
		return config, fmt.Errorf("couldn't read file: %w", err)
	}

	switch filepath.Ext(path) {
	case ".json":
		err := json.Unmarshal(b, &config)
		if err != nil {
			return config, fmt.Errorf("error reading file %s: %w", path, err)
		}
	case ".yml", ".yaml":
		err := yaml.Unmarshal(b, &config)
		if err != nil {
			return config, fmt.Errorf("error reading file %s: %w", path, err)
		}
	}

	config.setConfigFileField(path)

	return config, nil
}

func getFiles(path string) ([]string, error) {
	extensions := []string{".json", ".yaml", ".yml"}

	paths := []string{}

	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && slices.Contains(extensions, filepath.Ext(info.Name())) {
			paths = append(paths, path)
		}

		return nil
	})
	if err != nil {
		return paths, fmt.Errorf("couldn't get config files: %w", err)
	}

	return paths, nil
}
