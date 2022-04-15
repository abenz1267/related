package main

import (
	"log"
	"os"
	"path/filepath"

	"golang.org/x/exp/slices"
)

var validFragments = []string{"type", "group"} //nolint

func main() {
	args := os.Args[1:]
	cfg := parseConfig()

	validate(args, cfg)
	create(args, cfg)
}

func create(args []string, cfg Config) {
	if args[0] == validFragments[0] {
		createType(cfg, args[1], args[2])

		return
	}

	group := getFragment(cfg.Groups, args[1])

	for _, v := range group.Types {
		createType(cfg, v, args[2])
	}
}

func createType(cfg Config, typename, filename string) {
	t := getFragment(cfg.Types, typename)

	err := os.MkdirAll(t.Path, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	_, err = os.OpenFile(filepath.Join(t.Path, filename+"."+t.Extension), os.O_RDWR|os.O_CREATE|os.O_EXCL, 0o666)
	if err != nil {
		log.Fatal(err)
	}
}

func getFragment[T Fragment](fragments []T, name string) T {
	var res T

	for _, v := range fragments {
		if v.getName() == name {
			res = v
		}
	}

	return res
}

func validate(args []string, cfg Config) {
	if !slices.Contains(validFragments, args[0]) {
		log.Fatalf("'%s' is not a valid fragment", args[0])
	}

	if args[0] == validFragments[0] {
		if !containsName(cfg.Types, args[1]) {
			log.Fatalf("Type with name `%s` doesn't exist.", args[1])
		}
	}

	if args[0] == validFragments[1] {
		if !containsName(cfg.Groups, args[1]) {
			log.Fatalf("Group with name `%s` doesn't exist.", args[1])
		}
	}
}

func containsName[T Fragment](fragments []T, name string) bool {
	var contains bool

	for _, v := range fragments {
		if v.getName() == name {
			contains = true

			break
		}
	}

	return contains
}
