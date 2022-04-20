package creation

import (
	"log"
	"os"
	"path/filepath"

	. "github.com/abenz1267/gonerics" //nolint
	"github.com/abenz1267/related/config"
)

const (
	TypeCmd    = "type"
	GroupCmd   = "group"
	LuaExt     = ".lua"
	JsExt      = ".js"
	BinaryExt  = ""
	PermFolder = 0o755
	PermFile   = 0o644
)

type CmdArgs struct {
	Kind      string
	Component string
	Name      string
}

func Create(args CmdArgs) {
	cfg := config.ReadConfigs()

	if args.Kind == TypeCmd {
		if len(cfg.Types) == 0 {
			log.Println("<no types found>")

			return
		}

		createType(cfg, args.Component, args.Name)

		return
	}

	if len(cfg.Groups) == 0 {
		log.Println("<no groups found>")

		return
	}

	createGroup(cfg, args)
}

func createGroup(cfg config.Config, args CmdArgs) {
	group := config.GetFragment(cfg.Groups, args.Component)

	execScript(group.Pre, args.Name, group)

	for _, v := range group.Types {
		createType(cfg, v, args.Name)
	}

	execScript(group.Post, args.Name, group)
}

func createType(cfg config.Config, typename, name string) {
	defer RecoverPrint()

	fragment := config.GetFragment(cfg.Types, typename)

	execScript(fragment.Pre, name, fragment)

	writeFile(fragment, name)

	execScript(fragment.Post, name, fragment)
}

func writeFile(fragment config.Type, name string) {
	path := filepath.Join(fragment.Path, name+fragment.Suffix)

	Try(os.MkdirAll(filepath.Dir(path), PermFolder))

	file := TryResult(os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_EXCL, PermFile))
	defer file.Close()

	if fragment.Template != "" {
		buffer := getTemplateData(fragment, name)

		TryResult(file.Write(buffer.Bytes()))
	}

	log.Printf("created: %s\n", path)
}
