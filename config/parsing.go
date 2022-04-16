package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const (
	ConfigFile    = "related.json"
	DotConfigFile = ".related.json"
)

type Fragment interface {
	GetName() string
	GetPath() string
	GetExt() string
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

func (t Type) GetExt() string {
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

func (g Group) GetExt() string {
	return ""
}

type Config struct {
	Types  []Type  `json:"types"`
	Groups []Group `json:"groups"`
}

func ReadConfig() Config {
	workingDir, err := os.Getwd()
	if err != nil {
		log.Panic(err)
	}

	data, err := ioutil.ReadFile(filepath.Join(workingDir, DotConfigFile))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			data, err = ioutil.ReadFile(filepath.Join(workingDir, ConfigFile))
			if err != nil {
				log.Panic(err)
			}
		} else {
			log.Panic(err)
		}
	}

	var cfg Config

	err = json.Unmarshal(data, &cfg)
	if err != nil {
		log.Panic(err)
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
