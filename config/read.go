package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/exp/slices"
	"gopkg.in/yaml.v2"
)

type Fragment struct {
	Name     string `json:"name" yaml:"name"`
	Path     string `json:"path" yaml:"path"`
	Template string `json:"template" yaml:"template"`
	Pre      string `json:"pre" yaml:"pre"`
	Post     string `json:"post" yaml:"post"`
	Suffix   string `json:"suffix" yaml:"suffix"`
}

const Clear = "CLEAR"

func (fragment *Fragment) inheritFrom(parent Fragment) {
	fragment.Post = clearOrInherit(fragment.Post, parent.Post)
	fragment.Pre = clearOrInherit(fragment.Pre, parent.Pre)
	fragment.Template = clearOrInherit(fragment.Template, parent.Template)
	fragment.Suffix = clearOrInherit(fragment.Suffix, parent.Suffix)

	path := filepath.Join(parent.Path, fragment.Path)

	if strings.Contains(path, Clear) {
		paths := strings.Split(path, Clear)

		if len(paths) > 1 {
			fragment.Path = paths[1]
		} else {
			fragment.Path = ""
		}
	} else {
		fragment.Path = path
	}
}

func clearOrInherit(child, parent string) string {
	if child == "" {
		child = parent
	} else if child == Clear {
		child = ""
	}

	return child
}

type Group struct {
	Name      string   `json:"name" yaml:"name"`
	Fragments []string `json:"fragments" yaml:"fragments"`
}

type Config struct {
	Parents   []Fragment `json:"parents" yaml:"parents"`
	Fragments []Fragment `json:"fragments" yaml:"fragments"`
	Groups    []Group    `json:"groups" yaml:"groups"`
}

func (config *Config) addFragment(fragment Fragment) {
	config.Fragments = append(config.Fragments, fragment)
}

func (config *Config) addGroup(group Group) {
	config.Groups = append(config.Groups, group)
}

func (config *Config) addParent(parent Fragment) {
	config.Parents = append(config.Parents, parent)
}

func (config *Config) merge(other Config) {
	for _, v := range other.Parents {
		config.addParent(v)
	}

	for _, v := range other.Fragments {
		config.addFragment(v)
	}

	for _, v := range other.Groups {
		config.addGroup(v)
	}
}

func (config *Config) transform() (err error) {
	for i, v := range config.Fragments {
		if strings.Contains(v.Name, "/") {
			parent, err := config.getParent(strings.Split(v.Name, "/")[0])
			if err != nil {
				return err
			}

			config.Fragments[i].inheritFrom(parent)
		}
	}

	return nil
}

var ErrParentNotFound = errors.New("parent not found")

func (config Config) getParent(name string) (Fragment, error) {
	for _, v := range config.Parents {
		if v.Name == name {
			return v, nil
		}
	}

	return Fragment{}, ErrParentNotFound
}

var ErrFragmentNotFound = errors.New("fragment not found")

func (config Config) getFragment(name string) (Fragment, error) {
	for _, v := range config.Fragments {
		if v.Name == name {
			return v, nil
		}
	}

	return Fragment{}, ErrFragmentNotFound
}

func Get() (config Config, err error) {
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		return
	}

	folders := []string{filepath.Join(userConfigDir, "related"), ".related"}

	for _, v := range folders {
		if _, err := os.Stat(v); errors.Is(err, os.ErrNotExist) {
			continue
		}

		files, err := getFiles(v)
		if err != nil {
			return config, err
		}

		for _, f := range files {
			other, err := read(f)
			if err != nil {
				return config, err
			}

			config.merge(other)
		}
	}

	err = config.transform()

	return config, err
}

func read(path string) (config Config, err error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return
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

	return
}

func getFiles(path string) (paths []string, err error) {
	extensions := []string{".json", ".yaml", ".yml"}

	err = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && slices.Contains(extensions, filepath.Ext(info.Name())) {
			paths = append(paths, path)
		}

		return nil
	})
	if err != nil {
		return
	}

	return
}
