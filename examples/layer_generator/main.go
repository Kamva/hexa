package main

import (
	"path"

	"github.com/kamva/gutil"
	"github.com/kamva/hexa/lg"
)

type TemplateData struct {
	Package   string
	Name      string // struct name for the implementation of our interface
	Interface *lg.Interface
}

func main() {
	src := path.Join(gutil.SourcePath(), "app.go")
	tmpl := path.Join(gutil.SourcePath(), "err_layer.tmpl")
	output := path.Join(gutil.SourcePath(), "err_layer.go")

	pkg, err := lg.NewPackageFromFilenames("github.com/kamva/hexa/examples/layer_generator", src)
	if err != nil {
		panic(err)
	}

	resolver := lg.NewEmbeddedResolver(nil, pkg)
	if err := resolver.Resolve(); err != nil {
		panic(err)
	}

	_, iface := pkg.FindInterface("App")

	data := &TemplateData{
		Package:   "main",
		Name:      "errLayer",
		Interface: iface,
	}

	if err := lg.GenerateLayer(tmpl, lg.Funcs(), output, data, true); err != nil {
		panic(err)
	}
}
