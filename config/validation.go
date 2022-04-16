package config

import (
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

	return missingFiles, missingComponents
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
