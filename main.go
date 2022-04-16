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

	lua "github.com/yuin/gopher-lua"
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
	case "list":
		if len(args) == 1 {
			log.Println("Please provide a type, either 'scripts' or 'templates'")
		}

		if len(args) == 2 {
			listAvailables(localFS, args[1], "")

			return
		}

		if len(args) == 3 {
			listAvailables(localFS, args[1], args[2])

			return
		}
	default:
		validate(args, cfg)
		create(args, cfg, localFS)
	}
}

func listAvailables(localFS fs.FS, kind string, parent string) {
	systems := []fs.FS{embeddedFS, localFS}

	files := make(map[string][]string)

	root := filepath.Join(".", kind)
	if parent != "" {
		root = filepath.Join(root, parent)
	}

	suffix := ".tmpl"
	if kind == "scripts" {
		suffix = ".lua"
	}

	for _, v := range systems {
		err := fs.WalkDir(v, root, func(path string, entry fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if !entry.IsDir() { //nolint
				if strings.HasSuffix(entry.Name(), suffix) {
					parent := strings.Split(path, string(filepath.Separator))[1]

					if _, ok := files[parent]; !ok {
						files[parent] = []string{}
					}

					tmplName := strings.TrimSuffix(entry.Name(), suffix)

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

	if len(files) == 0 {
		fmt.Println("<nothing found>")

		return
	}

	for k, v := range files {
		fmt.Printf("%s: %s\n", k, strings.Join(v, ",")) //nolint
	}
}

func execScript(scriptFile string, filename string, localFS fs.FS) {
	script, err := fs.ReadFile(localFS, filepath.Join("scripts", scriptFile+".lua"))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			log.Printf("Script not found: %s.", scriptFile)
		}

		return
	}

	L := lua.NewState()
	defer L.Close()

	L.SetGlobal("Name", lua.LString(filename))

	err = L.DoString(string(script))
	if err != nil {
		log.Panic(err)
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

	if group.Pre != "" {
		execScript(group.Pre, args[2], localFS)
	}

	for _, v := range group.Types {
		createType(cfg, v, args[2], localFS)
	}

	if group.Post != "" {
		execScript(group.Post, args[2], localFS)
	}
}

func createType(cfg Config, typename, filename string, localFS fs.FS) {
	fragment := getFragment(cfg.Types, typename)

	if fragment.Pre != "" {
		execScript(fragment.Pre, filename, localFS)
	}

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

	if fragment.Post != "" {
		execScript(fragment.Post, filename, localFS)
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
