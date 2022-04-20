package creation

import (
	"bytes"
	"os"
	"path/filepath"
	"text/template"

	. "github.com/abenz1267/gonerics" //nolint
	"github.com/abenz1267/related/config"
	"github.com/abenz1267/related/files"
)

func getTemplateData(fragment config.Type, name string) bytes.Buffer {
	workingDir := TryResult(os.Getwd())

	fullPath := filepath.Join(fragment.GetPath(), name)

	path := filepath.Dir(fullPath)
	if path == "." {
		path = ""
	}

	data := map[string]string{
		"workingDir": workingDir,
		"path":       path,
		"name":       filepath.Base(fullPath),
		"suffix":     fragment.GetSuffix(),
	}

	var buffer bytes.Buffer

	path, system := files.FindFile(fragment.Template, files.TemplateDir)

	tmpl := TryResult(template.ParseFS(system, path))

	Try(tmpl.Execute(&buffer, data))

	return buffer
}
