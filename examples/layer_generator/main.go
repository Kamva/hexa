package main

import (
	"fmt"
	"path"

	"github.com/kamva/gutil"
	"github.com/kamva/hexa/lg"
	"github.com/kamva/tracer"
)

type TemplateData struct {
	Package   string
	Name      string // struct name for the implementation of our interface
	Interface *lg.InterfaceMetadata
}

func main() {
	src := path.Join(gutil.SourcePath(), "app.go")
	tmpl := path.Join(gutil.SourcePath(), "err_layer.tmpl")
	output := path.Join(gutil.SourcePath(), "err_layer.go")

	metadata, err := lg.ExtractInterfaceMetadata(src, "App")
	if err != nil {
		fmt.Println(tracer.StackAsString(tracer.Trace(err)))
		panic(err)
	}

	data := &TemplateData{
		Package:   "main",
		Name:      "errLayer",
		Interface: metadata,
	}

	if err := lg.GenerateLayer(tmpl, lg.Funcs(), output, data); err != nil {
		panic(err)
	}
}
