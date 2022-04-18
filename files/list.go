package files

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/exp/slices"
)

const (
	ListCmd        = "list"
	ListInvalid    = 1
	ListWithType   = 2
	ListWithParent = 3
)

func List(args []string) {
	switch len(args) {
	case ListInvalid:
		log.Printf("Please provide a type, either '%s' or '%s'", ScriptDir, TemplateDir)
	case ListWithType:
		if !slices.Contains([]TypeDir{ScriptDir, TemplateDir}, TypeDir(args[1])) {
			log.Fatalf("Invalid collection '%s'", args[1])
		}

		display(Availables(args[1], ""))
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

func Availables(kind string, parent string) map[string][]string {
	systems := Systems()

	result := make(map[string][]string)

	root := filepath.Join(".", kind)
	if parent != "" {
		root = filepath.Join(root, parent)
	}

	for _, v := range systems {
		if v == nil {
			continue
		}

		err := fs.WalkDir(v, root, func(path string, entry fs.DirEntry, err error) error {
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
		})
		if err != nil {
			if !errors.Is(err, os.ErrNotExist) {
				log.Fatal(err)
			}
		}
	}

	if len(result) == 0 {
		return nil
	}

	return result
}
