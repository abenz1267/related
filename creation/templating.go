package creation

import (
	"bytes"
	"path/filepath"
	"text/template"

	. "github.com/abenz1267/gonerics" //nolint
	"github.com/abenz1267/related/files"
)

func getTemplateData(templateName, name string) bytes.Buffer {
	var buffer bytes.Buffer

	path, system := files.FindFile(templateName, files.TemplateDir)

	tmpl := TryResult(template.ParseFS(system, path))

	Try(tmpl.Execute(&buffer, filepath.Base(name)))

	return buffer
}
