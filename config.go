package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type Fragment interface {
	getName() string
}

type Type struct {
	Name      string `json:"name"`
	Path      string `json:"path"`
	Template  string `json:"template"`
	Pre       string `json:"pre"`
	Post      string `json:"post"`
	Extension string `json:"extension"`
}

func (t Type) getName() string {
	return t.Name
}

type Group struct {
	Name    string   `json:"name"`
	Types   []string `json:"types"`
	Pre       string `json:"pre"`
	Post      string `json:"post"`
}

func (g Group) getName() string {
	return g.Name
}

type Config struct {
	Types  []Type  `json:"types"`
	Groups []Group `json:"groups"`
}

func parseConfig() Config {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	cfgPath := filepath.Join(wd, ".related.json")

	data, err := ioutil.ReadFile(cfgPath)
	if err != nil {
		log.Fatal(err)
	}

	var cfg Config
	json.Unmarshal(data, &cfg)

	return cfg
}
