package config

import (
	"fmt"
	"io/fs"
	"log"
	"path/filepath"
	"strings"

	. "github.com/abenz1267/gonerics" //nolint
	"github.com/abenz1267/related/files"
	"golang.org/x/exp/slices"
)

const (
	ListCmd        = "list"
	ListInvalid    = 1
	ListWithType   = 2
	ListWithParent = 3
)

func List(args []string) {
	types := Slice("templates", "scripts", "groups", "types")

	switch len(args) {
	case ListInvalid:
		log.Printf("Please provide a type, either '%s' or '%s'", files.ScriptDir, files.TemplateDir)
	case ListWithType:
		if !slices.Contains(types, args[1]) {
			log.Fatalf("Invalid collection '%s'", args[1])
		}

		switch args[1] {
		case types[0], types[1]:
			display(Availables(args[1], ""))
		case types[2]:
			display(fragments(ReadConfigs().Groups, types[2]))
		case types[3]:
			display(fragments(ReadConfigs().Types, types[3]))
		}

	case ListWithParent:
		display(Availables(args[1], args[2]))
	}
}

func display(data map[string][]string) {
	if len(data) == 0 {
		log.Println("<nothing found>")

		return
	}

	for k, v := range data {
		fmt.Printf("%s:\n", k)

		for i := 0; i <= len(k); i++ {
			fmt.Print("-")
		}
		fmt.Print("\n")

		for _, n := range v {
			fmt.Println(n)
		}

		fmt.Print("\n\n")
	}
}

func fragments[T Fragment](fragments []T, parent string) map[string][]string {
	result := map[string][]string{}

	for _, v := range fragments {
		result[parent] = append(result[parent], v.GetName())
	}

	return result
}

func Availables(kind string, parent string) map[string][]string {
	defer RecoverFatal()

	systems := files.Systems()

	result := map[string][]string{}

	root := filepath.Join(".", kind)
	if parent != "" {
		root = filepath.Join(root, parent)
	}

	for _, v := range systems {
		if v == nil {
			continue
		}

		Try(fs.WalkDir(v, root, func(path string, entry fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if !entry.IsDir() {
				if parent != "" {
					parts := strings.Split(path, "/")
					result[parent] = append(result[parts[1]], strings.Join(parts[2:], "/"))
				} else {
					parts := strings.Split(path, "/")
					result[parts[1]] = append(result[parts[1]], strings.Join(parts[2:], "/"))
				}
			}

			return nil
		}))
	}

	if len(result) == 0 {
		return nil
	}

	return result
}
