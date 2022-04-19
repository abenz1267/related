package config

import (
	"log"
	"os"
	"strings"

	"github.com/abenz1267/related/files"
	"golang.org/x/exp/slices"
)

func Validate(cfg Config) ([]string, []string) {
	missingFiles := []string{}
	missingComponents := []string{}
	names := []string{}

	for _, component := range cfg.Types {
		names = append(names, component.Name)

		if !exists(component.Template, files.TemplateDir) {
			missingFiles = append(missingFiles, component.Template)
		}

		if !exists(component.Pre, files.ScriptDir) {
			missingFiles = append(missingFiles, component.Pre)
		}

		if !exists(component.Post, files.ScriptDir) {
			missingFiles = append(missingFiles, component.Post)
		}
	}

	for _, group := range cfg.Groups {
		for _, name := range group.Types {
			if !slices.Contains(names, name) {
				missingComponents = append(missingComponents, name)
			}
		}

		if !exists(group.Pre, files.ScriptDir) {
			missingFiles = append(missingFiles, group.Pre)
		}

		if !exists(group.Post, files.ScriptDir) {
			missingFiles = append(missingFiles, group.Post)
		}
	}

	var abort bool

	ambiguousTypes := findAmbiguous(cfg.Types)

	if len(ambiguousTypes) > 0 {
		log.Printf("Ambiguous types: %s", strings.Join(ambiguousTypes, ", "))

		abort = true
	}

	ambiguousGroups := findAmbiguous(cfg.Groups)

	if len(ambiguousGroups) > 0 {
		log.Printf("Ambiguous groups: %s", strings.Join(ambiguousGroups, ", "))

		abort = true
	}

	if abort {
		os.Exit(1)
	}

	return missingFiles, missingComponents
}

func findAmbiguous[T Fragment](fragments []T) []string {
	exist := map[string]struct{}{}
	ambigious := []string{}

	for _, v := range fragments {
		_, exists := exist[v.GetName()]
		if !exists {
			exist[v.GetName()] = struct{}{}
		} else {
			ambigious = append(ambigious, v.GetName())
		}
	}

	return ambigious
}

func exists(name string, dir files.TypeDir) bool {
	if name != "" {
		if file, _ := files.FindFile(name, dir); file == "" {
			return false
		}

		return true
	}

	return true
}
