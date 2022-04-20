package creation

import (
	"errors"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	. "github.com/abenz1267/gonerics" //nolint
	"github.com/abenz1267/related/config"
	"github.com/abenz1267/related/files"
	lua "github.com/yuin/gopher-lua"
)

const (
	LuaArgWorkingDir = 1
	LuaArgPath       = 2
	LuaArgName       = 3
	LuaArgSuffix     = 4
)

func execLua(system fs.FS, fragment config.Fragment, path, name string) {
	file := TryResult(fs.ReadFile(system, path))

	state := lua.NewState()
	defer state.Close()

	workingDir := TryResult(os.Getwd())

	fullPath := filepath.Join(fragment.GetPath(), name)

	args := lua.LTable{
		Metatable: nil,
	}
	args.Insert(LuaArgWorkingDir, lua.LString(workingDir))
	args.Insert(LuaArgPath, lua.LString(filepath.Dir(fullPath)))
	args.Insert(LuaArgName, lua.LString(filepath.Base(fullPath)))
	args.Insert(LuaArgSuffix, lua.LString(fragment.GetSuffix()))

	state.SetGlobal("arg", &args)

	Try(state.DoString(string(file)))
}

func execBinaryOrJavascript(fragment config.Fragment, path, name string) {
	wd := TryResult(os.Getwd())

	fullPath := filepath.Join(fragment.GetPath(), name)

	args := Slice(path, wd, filepath.Dir(fullPath), filepath.Base(fullPath), fragment.GetSuffix())

	projectPath := filepath.Join(string(files.ProjectDir), path)

	_, err := os.Stat(projectPath)
	if !errors.Is(err, os.ErrNotExist) {
		args[0] = projectPath
	} else {
		userConfigDir := TryResult(os.UserConfigDir())

		args[0] = filepath.Join(userConfigDir, string(files.ConfigDir), path)
	}

	if filepath.Ext(path) == JsExt {
		args = append(Slice("node"), args...)
	}

	cmd := exec.Command(args[0], args[1:]...)

	if res := TryResult(cmd.Output()); len(res) != 0 {
		log.Println(string(res))
	}
}

func execScript(script string, name string, fragment config.Fragment) {
	if script == "" {
		return
	}

	path, system := files.FindFile(script, files.ScriptDir)

	log.Printf("executing: %s", script)

	switch filepath.Ext(script) {
	case LuaExt:
		execLua(system, fragment, path, name)
	case JsExt, BinaryExt:
		execBinaryOrJavascript(fragment, path, name)
	default:
		log.Fatalf("Unknown script extension '%s'", filepath.Ext(script))
	}
}
