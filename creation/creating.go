package creation

import (
	"log"
	"os"
	"path/filepath"

	"github.com/abenz1267/related/config"
)

const (
	TypeCmd   = "type"
	GroupCmd  = "group"
	LuaExt    = ".lua"
	JsExt     = ".js"
	BinaryExt = ""
)

type CmdArgs struct {
	Kind      string
	Component string
	Name      string
}

func Create(args CmdArgs) {
	cfg := config.ReadConfig()

	if args.Kind == TypeCmd {
		createType(cfg, args.Component, args.Name)

		return
	}

	createGroup(cfg, args)
}

func createGroup(cfg config.Config, args CmdArgs) {
	group := getFragment(cfg.Groups, args.Component)

	execScript(group.Pre, args.Name, group)

	for _, v := range group.Types {
		createType(cfg, v, args.Name)
	}

	execScript(group.Post, args.Name, group)
}

func createType(cfg config.Config, typename, name string) {
	fragment := getFragment(cfg.Types, typename)

	execScript(fragment.Pre, name, fragment)

	path := filepath.Join(fragment.Path, name+fragment.Suffix)

	err := os.MkdirAll(filepath.Dir(path), 0o755)
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	if fragment.Template != "" {
		buffer := getTemplateData(fragment.Template, name)

		_, err := file.Write(buffer.Bytes())
		if err != nil {
			log.Panic(err)
		}
	}

	execScript(fragment.Post, name, fragment)
}

func getFragment[T config.Fragment](fragments []T, name string) T { //nolint
	var res T

	for _, v := range fragments {
		if v.GetName() == name {
			res = v
		}
	}

	return res
}
