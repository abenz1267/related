package config

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"golang.org/x/exp/slices"
)

const Clear = "<clear>"

type Fragment struct {
	Name       string `json:"name" yaml:"name"`
	Template   string `json:"template" yaml:"template"`
	Script     string `json:"script" yaml:"script"`
	File       string `json:"file" yaml:"file"`
	ConfigFile string
}

func (fragment Fragment) getParentName() string {
	if !strings.Contains(fragment.Name, "/") {
		return ""
	}

	return strings.Split(fragment.Name, "/")[0]
}

func (fragment *Fragment) inheritFrom(parent Fragment) {
	fragment.Script = clearOrInherit(fragment.Script, parent.Script)
	fragment.Template = clearOrInherit(fragment.Template, parent.Template)
	fragment.File = clearOrInherit(fragment.File, parent.File)

	if strings.Contains(fragment.File, Clear) {
		fragment.File = strings.TrimPrefix(fragment.File, Clear)
		fragment.File = strings.TrimPrefix(fragment.File, "/")
	} else {
		fragment.File = filepath.Join(parent.File, fragment.File)
	}
}

func clearOrInherit(child, parent string) string {
	if child == Clear {
		return ""
	}

	return parent
}

type Group struct {
	Name       string `json:"name" yaml:"name"`
	ConfigFile string
	Script     string   `json:"script" yaml:"script"`
	Fragments  []string `json:"fragments" yaml:"fragments"`
}

type Config struct {
	Parents   []Fragment `json:"parents" yaml:"parents"`
	Fragments []Fragment `json:"fragments" yaml:"fragments"`
	Groups    []Group    `json:"groups" yaml:"groups"`
}

func (config *Config) setConfigFileField(path string) {
	for i := range config.Fragments {
		config.Fragments[i].ConfigFile = path
	}

	for i := range config.Parents {
		config.Parents[i].ConfigFile = path
	}

	for i := range config.Groups {
		config.Groups[i].ConfigFile = path
	}
}

func (config Config) getFragmentNames() []string {
	res := []string{}

	for _, v := range config.Fragments {
		res = append(res, v.Name)
	}

	return res
}

func (config Config) getParentNames() []string {
	res := []string{}

	for _, v := range config.Parents {
		res = append(res, v.Name)
	}

	return res
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

func (config *Config) transform() error {
	for i, v := range config.Fragments {
		if strings.Contains(v.Name, "/") {
			parent, err := config.GetFragment(strings.Split(v.Name, "/")[0], config.Parents)
			if err != nil {
				return fmt.Errorf("parent not found: %w", err)
			}

			config.Fragments[i].inheritFrom(parent)
		}
	}

	return nil
}

var (
	ErrFragmentNotFound = errors.New("fragment not found")
	ErrGroupNotFound    = errors.New("group not found")
)

func (config Config) GetFragment(name string, list []Fragment) (Fragment, error) {
	for _, v := range list {
		if v.Name == name {
			return v, nil
		}
	}

	return Fragment{}, ErrFragmentNotFound
}

func (config Config) GetGroup(name string, list []Group) (Group, error) {
	for _, v := range list {
		if v.Name == name {
			return v, nil
		}
	}

	return Group{}, ErrGroupNotFound
}

func (config *Config) sort() {
	slices.SortFunc(config.Parents, func(a, b Fragment) bool {
		return a.Name < b.Name
	})

	slices.SortFunc(config.Fragments, func(a, b Fragment) bool {
		return a.Name < b.Name
	})

	slices.SortFunc(config.Groups, func(a, b Group) bool {
		return a.Name < b.Name
	})
}
