package creation

import (
	"bytes"
	"log"
	"path/filepath"
	"text/template"

	"github.com/abenz1267/related/files"
)

func getTemplateData(templateName, name string) bytes.Buffer {
	var buffer bytes.Buffer

	path, system := files.FindFile(templateName, files.TemplateDir)

	tmpl, err := template.ParseFS(system, path)
	if err != nil {
		log.Fatal(err)
	}

	err = tmpl.Execute(&buffer, filepath.Base(name))
	if err != nil {
		log.Fatal(err)
	}

	return buffer
}
