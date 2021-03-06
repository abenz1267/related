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

	. "github.com/abenz1267/gonerics" //nolint
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
	defer RecoverFatal()

	definitions := Slice(DotConfigFile, ConfigFile)

	Try(filepath.Walk(string(files.ProjectDir), func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && slices.Contains(validExtensions, filepath.Ext(info.Name())) {
			definitions = append(definitions, path)
		}

		return nil
	}))

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
		log.Panic("No definitions found.")
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

	return transformParents(cfg)
}

func transformParents(cfg Config) Config {
	for i, v := range cfg.Types {
		if strings.Contains(v.Name, "/") {
			types := strings.Split(v.Name, "/")
			parent := GetFragment(cfg.Types, types[0])
			cfg.Types[i] = merge(parent, v)
			cfg.Types[i].Name = v.Name
		}
	}

	return cfg
}

const clear = "CLEAR"

func merge(parent, child Type) Type {
	if child.Suffix != "" {
		parent.Suffix = child.Suffix
	}

	if child.Suffix == clear {
		parent.Suffix = ""
	}

	if child.Path != "" {
		parent.Path = filepath.Join(parent.Path, child.Path)
	}

	if child.Path == clear {
		parent.Path = ""
	}

	if child.Pre != "" {
		parent.Pre = child.Pre
	}

	if child.Pre == clear {
		parent.Pre = ""
	}

	if child.Post != "" {
		parent.Post = child.Post
	}

	if child.Post == clear {
		parent.Post = ""
	}

	if child.Template != "" {
		parent.Template = child.Template
	}

	if child.Template == clear {
		parent.Template = ""
	}

	return parent
}

func readConfig(path string) Config {
	defer RecoverPanic()

	data, err := ioutil.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return Config{} //nolint
		}

		log.Panic(err)
	}

	var cfg Config

	if strings.HasSuffix(path, validExtensions[0]) {
		Try(json.Unmarshal(data, &cfg))
	} else {
		Try(yaml.Unmarshal(data, &cfg))
	}

	return cfg
}

func GetFragment[T Fragment](fragments []T, name string) T { //nolint
	var res T

	for _, v := range fragments {
		if v.GetName() == name {
			res = v
		}
	}

	return res
}
