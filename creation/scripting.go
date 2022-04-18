package creation

import (
	"errors"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/abenz1267/related/config"
	"github.com/abenz1267/related/files"
	lua "github.com/yuin/gopher-lua"
)

func execLua(system fs.FS, fragment config.Fragment, path, name string) {
	file, err := fs.ReadFile(system, path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			log.Printf("Script not found: %s.", filepath.Base(path))
		}

		return
	}

	state := lua.NewState()
	defer state.Close()

	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	fullPath := filepath.Join(fragment.GetPath(), name)

	args := lua.LTable{
		Metatable: nil,
	}
	args.Insert(1, lua.LString(wd))
	args.Insert(2, lua.LString(filepath.Dir(fullPath)))
	args.Insert(3, lua.LString(filepath.Base(fullPath)))
	args.Insert(4, lua.LString(fragment.GetSuffix()))

	state.SetGlobal("arg", &args)

	err = state.DoString(string(file))
	if err != nil {
		log.Panic(err)
	}
}

func execBinaryOrJavascript(fragment config.Fragment, path, name string) {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	fullPath := filepath.Join(fragment.GetPath(), name)

	args := []string{path, wd, filepath.Dir(fullPath), filepath.Base(fullPath), fragment.GetSuffix()}

	projectPath := filepath.Join(string(files.ProjectDir), path)

	_, err = os.Stat(projectPath)
	if !errors.Is(err, os.ErrNotExist) {
		args[0] = projectPath
	} else {
		wd, wdErr := os.UserConfigDir()
		if wdErr != nil {
			log.Fatal(err)
		}

		args[0] = filepath.Join(wd, string(files.ConfigDir), path)
	}

	if filepath.Ext(path) == JsExt {
		args = append([]string{"node"}, args...)
	}

	cmd := exec.Command(args[0], args[1:]...)

	res, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}

	if len(res) != 0 {
		log.Println(string(res))
	}
}

func execScript(script string, name string, fragment config.Fragment) {
	if script == "" {
		return
	}

	path, system := files.FindFile(script, files.ScriptDir)

	switch filepath.Ext(script) {
	case LuaExt:
		execLua(system, fragment, path, name)
	case JsExt, BinaryExt:
		execBinaryOrJavascript(fragment, path, name)
	default:
		log.Fatalf("Unknown script extension '%s'", filepath.Ext(script))
	}
}
