package testingcommons

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/abenz1267/related/config"
	"github.com/abenz1267/related/files"
)

const (
	ConfigDir       = "config"
	ConfigRelated   = "related"
	Parent          = "parent"
	ProjectTemplate = "template.tmpl"
	ConfigTemplate  = "configtemplate.tmpl"
	ProjectScript   = "script.lua"
	ConfigScript    = "configscript.lua"
)

const (
	ProjectFolder = string(files.ProjectDir)
	TemplateDir   = string(files.TemplateDir)
	ScriptDir     = string(files.ScriptDir)
)

func GetName(file string) string {
	return strings.Join([]string{Parent, file}, "/")
}

func CreateTmpData() {
	files.SetConfigDir(ConfigDir)

	testFolders := [][]string{
		{ProjectFolder, TemplateDir, Parent},
		{ProjectFolder, ScriptDir, Parent},
		{ConfigDir, ConfigRelated, TemplateDir, Parent},
		{ConfigDir, ConfigRelated, ScriptDir, Parent},
	}

	testFiles := [][]string{
		{ProjectFolder, TemplateDir, Parent, ProjectTemplate},
		{ProjectFolder, ScriptDir, Parent, ProjectScript},
		{ConfigDir, ConfigRelated, TemplateDir, Parent, ConfigTemplate},
		{ConfigDir, ConfigRelated, ScriptDir, Parent, ConfigScript},
	}

	for _, v := range testFolders {
		err := os.MkdirAll(filepath.Join(v...), os.ModePerm)
		if err != nil {
			panic(err)
		}
	}

	for _, v := range testFiles {
		if _, err := os.Create(filepath.Join(v...)); err != nil {
			panic(err)
		}
	}
}

func Cleanup() {
	if err := os.RemoveAll(string(files.ProjectDir)); err != nil {
		panic(err)
	}

	if err := os.RemoveAll(ConfigDir); err != nil {
		panic(err)
	}

	if err := os.RemoveAll(config.ConfigFile); err != nil {
		panic(err)
	}
}
