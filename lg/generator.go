package lg

import (
	"bytes"
	"io/ioutil"
	"path"
	"text/template"

	"github.com/kamva/tracer"
	"golang.org/x/tools/imports"
)

type TemplateData struct {
	Package   string
	Name      string // struct name for the implementation of our interface
	Interface *InterfaceMetadata
}

func GenerateLayer(tmpl string, funcs template.FuncMap, outputFile string, data any, reformat bool) error {
	t := template.Must(template.New(path.Base(tmpl)).Funcs(funcs).ParseFiles(tmpl))

	out := bytes.NewBufferString("")
	err := t.Execute(out, data)
	if err != nil {
		return tracer.Trace(err)
	}
	code := out.Bytes()

	if reformat {
		code, err = imports.Process(outputFile, code, &imports.Options{Comments: true})
		if err != nil {
			return tracer.Trace(err)
		}
	}

	return tracer.Trace(ioutil.WriteFile(outputFile, code, 0644))
}
