package creation

import (
	"bytes"
	"log"
	"text/template"

	"github.com/abenz1267/related/files"
)

func getTemplateData(templateName, name string) bytes.Buffer {
	var buffer bytes.Buffer

	filepath, system := files.FindFile(templateName, files.TemplateDir)

	tmpl, err := template.ParseFS(system, filepath)
	if err != nil {
		log.Fatal(err)
	}

	err = tmpl.Execute(&buffer, name)
	if err != nil {
		log.Fatal(err)
	}

	return buffer
}
