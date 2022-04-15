package main

import (
	"bytes"
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"golang.org/x/exp/slices"
)

var validFragments = []string{"type", "group"} //nolint

//go:embed templates
var embeddedFS embed.FS

func main() {
	cfgDir, err := os.UserConfigDir()
	if err != nil {
		log.Fatal(err)
	}

	localFS := os.DirFS(filepath.Join(cfgDir, "related"))

	args := os.Args[1:]
	cfg := parseConfig()

	switch args[0] {
	case "listtemplates":
		if len(args) == 2 {
			listTemplates(localFS, args[1])

			return
		}

		listTemplates(localFS, "")
	default:
		validate(args, cfg)
		create(args, cfg, localFS)
	}
}

func listTemplates(localFS fs.FS, parent string) {
	systems := []fs.FS{embeddedFS, localFS}

	files := make(map[string][]string)

	root := "."
	if parent != "" {
		root = filepath.Join("templates", parent)
	}

	for _, v := range systems {
		err := fs.WalkDir(v, root, func(path string, entry fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if !entry.IsDir() { //nolint
				if strings.HasSuffix(entry.Name(), ".tmpl") {
					parent := strings.Split(path, string(filepath.Separator))[1]

					if _, ok := files[parent]; !ok {
						files[parent] = []string{}
					}

					tmplName := strings.TrimSuffix(entry.Name(), ".tmpl")

					if !slices.Contains(files[parent], tmplName) {
						files[parent] = append(files[parent], tmplName)
					}
				}
			}

			return nil
		})
		if err != nil {
			if !errors.Is(err, os.ErrNotExist) {
				log.Fatal(err)
			}
		}
	}

	for k, v := range files {
		fmt.Printf("%s: %s\n", k, strings.Join(v, ","))
	}
}

func getTemplateData(name string, filename string, filesystem fs.FS) bytes.Buffer {
	templatename := strings.Split(name, "/")[1]

	var buffer bytes.Buffer

	err := fs.WalkDir(filesystem, ".", func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !entry.IsDir() { //nolint
			if strings.HasPrefix(entry.Name(), templatename) {
				tmpl, err := template.ParseFS(filesystem, path)
				if err != nil {
					log.Fatal(err)
				}

				err = tmpl.Execute(&buffer, filename)
				if err != nil {
					log.Fatal(err)
				}
			}
		}

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	return buffer
}

func create(args []string, cfg Config, localFS fs.FS) {
	if args[0] == validFragments[0] {
		createType(cfg, args[1], args[2], localFS)

		return
	}

	group := getFragment(cfg.Groups, args[1])

	for _, v := range group.Types {
		createType(cfg, v, args[2], localFS)
	}
}

func createType(cfg Config, typename, filename string, localFS fs.FS) {
	fragment := getFragment(cfg.Types, typename)

	err := os.MkdirAll(fragment.Path, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	filepath := filepath.Join(fragment.Path, filename+"."+fragment.Extension)

	file, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE|os.O_EXCL, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	if fragment.Template != "" {
		var buffer bytes.Buffer

		buffer = getTemplateData(fragment.Template, filename, localFS)
		if buffer.Len() == 0 {
			buffer = getTemplateData(fragment.Template, filename, embeddedFS)
		}

		if buffer.Len() == 0 {
			return
		}

		_, err := file.Write(buffer.Bytes())
		if err != nil {
			log.Panic(err)
		}
	}
}

func getFragment[T Fragment](fragments []T, name string) T { //nolint
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
