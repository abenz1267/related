package config

import (
	"encoding/json"
	"errors"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/abenz1267/related/files"
	"golang.org/x/exp/slices"
	"gopkg.in/yaml.v3"
)

const (
	ConfigFile    = "related.json"
	DotConfigFile = ".related.json"
)

type Fragment interface {
	GetName() string
	GetPath() string
	GetSuffix() string
}

type Type struct {
	Name     string `json:"name"`
	Path     string `json:"path"`
	Template string `json:"template"`
	Pre      string `json:"pre"`
	Post     string `json:"post"`
	Suffix   string `json:"suffix"`
}

func (t Type) GetName() string {
	return t.Name
}

func (t Type) GetPath() string {
	return t.Path
}

func (t Type) GetSuffix() string {
	return t.Suffix
}

type Group struct {
	Name  string   `json:"name"`
	Pre   string   `json:"pre"`
	Post  string   `json:"post"`
	Types []string `json:"types"`
}

func (g Group) GetName() string {
	return g.Name
}

func (g Group) GetPath() string {
	return ""
}

func (g Group) GetSuffix() string {
	return ""
}

type Config struct {
	Types  []Type  `json:"types"`
	Groups []Group `json:"groups"`
}

var validExtensions = []string{".json", ".yaml"}

func ReadConfigs() Config {
	definitions := []string{DotConfigFile, ConfigFile}

	err := filepath.Walk(string(files.ProjectDir), func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && slices.Contains(validExtensions, filepath.Ext(info.Name())) {
			definitions = append(definitions, path)
		}

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	cfg := Config{
		Types:  []Type{},
		Groups: []Group{},
	}

	for _, v := range definitions {
		definition := readConfig(v)

		if definition.Groups != nil {
			cfg.Groups = append(cfg.Groups, definition.Groups...)
		}

		if definition.Types != nil {
			cfg.Types = append(cfg.Types, definition.Types...)
		}
	}

	if len(cfg.Types) == 0 && len(cfg.Groups) == 0 {
		log.Fatal("No definitions found.")
	}

	var abort bool

	missingFiles, missingComponents := Validate(cfg)

	if len(missingFiles) > 0 {
		log.Printf("Missing files: %s", strings.Join(missingFiles, ", "))

		abort = true
	}

	if len(missingComponents) > 0 {
		log.Printf("Missing components: %s", strings.Join(missingComponents, ", "))

		abort = true
	}

	if abort {
		os.Exit(1)
	}

	return cfg
}

func readConfig(path string) Config {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return Config{} //nolint
		}

		log.Panic(err)
	}

	var cfg Config

	if strings.HasSuffix(path, validExtensions[0]) {
		err = json.Unmarshal(data, &cfg)
	} else {
		err = yaml.Unmarshal(data, &cfg)
	}

	if err != nil {
		log.Panic(err)
	}

	return cfg
}
