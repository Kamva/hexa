package lg

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"path"
	"reflect"

	"github.com/kamva/tracer"
)

const DefaultParseMode = parser.AllErrors | parser.ParseComments

type Import struct {
	Name string // For regular imports this value is empty. For alias imports it's the alias.
	Path string
}

type Package struct {
	Name  string  // e.g., hexa
	Path  string  // e.g., github.com/kamva/hexa
	Files []*File // The golang files
}

type File struct {
	ImportMap map[string]string // map's key is the imported package's name or its alias, value is the package's path. e.g., hexa -> github.com/kamva/hexa

	PackageName string // The package's name. e.g., hexa
	Imports     []*Import

	Interfaces []*Interface // map's key is the interface's name.
	Structs    []*Struct    // map's key is the struct name.
}

type Interface struct {
	Doc         string
	Annotations Annotations
	Name        string
	Embedded    []*EmbeddedField
	Methods     []*Method
}

type Struct struct {
	Doc         string
	Annotations Annotations
	Name        string
	Embedded    []*EmbeddedField
	Fields      []*Field
}

type Method struct {
	Doc         string
	Annotations Annotations
	Name        string
	Params      []*MethodParam
	Results     []*MethodResult
}

type MethodParam struct {
	Name string
	Type string
}

type MethodResult = MethodParam

type Field struct {
	Doc         string
	Annotations Annotations
	Name        string
	Type        string
	Tag         reflect.StructTag
}
type EmbeddedField struct {
	IsResolved  bool // When we add all fields of the embedded type to its parent type, it's resolved.
	Doc         string
	Annotations Annotations
	Type        string
	Tag         reflect.StructTag
}

func (p *Package) FindInterface(name string) (*File, *Interface) {
	for _, f := range p.Files {
		if iface := f.FindInterface(name); iface != nil {
			return f, iface
		}
	}
	return nil, nil
}
func (p *Package) FindStruct(name string) (*File, *Struct) {
	for _, f := range p.Files {
		if strct := f.FindStruct(name); strct != nil {
			return f, strct
		}
	}
	return nil, nil
}

func (f *File) FindInterface(name string) *Interface {
	for _, iface := range f.Interfaces {
		if iface.Name == name {
			return iface
		}
	}
	return nil
}

func (f *File) FindStruct(name string) *Struct {
	for _, strct := range f.Structs {
		if strct.Name == name {
			return strct
		}
	}
	return nil
}

func (i *Interface) MethodByName(name string) *Method {
	for _, m := range i.Methods {
		if m.Name == name {
			return m
		}
	}
	return nil
}

func (r MethodResult) joinNameAndType() string {
	if r.Name == "" {
		return r.Type
	}
	return fmt.Sprintf("%s %s", r.Name, r.Type)
}

func NewPackage(pkgPath string, files []*File) *Package {
	return &Package{
		Name:  path.Base(pkgPath),
		Path:  pkgPath,
		Files: files,
	}
}

func NewPackageFromAstPackage(pkgPath string, astPkg *ast.Package) *Package {
	files := make([]*File, len(astPkg.Files))
	var i int
	for _, f := range astPkg.Files {
		files[i] = NewFile(f)
		i++
	}
	return NewPackage(pkgPath, files)
}

func SinglePackageFromDir(pkgPath string, dir string) (*Package, error) {
	return SinglePackageFromDirWithOpts(pkgPath, dir, nil, DefaultParseMode)
}

func SinglePackageFromDirWithOpts(pkgPath string, dir string, filter func(fs.FileInfo) bool, mode parser.Mode) (*Package, error) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dir, filter, mode)
	if err != nil {
		return nil, tracer.Trace(err)
	}

	pkg, ok := pkgs[path.Base(pkgPath)]
	if !ok {
		return nil, tracer.Trace(errors.New("package not found"))
	}

	return NewPackageFromAstPackage(pkgPath, pkg), nil
}

func NewPackageFromFilenames(pkgPath string, filenames ...string) (*Package, error) {
	return NewPackageFromFilenamesWithOpts(pkgPath, DefaultParseMode, filenames...)
}

func NewPackageFromFilenamesWithOpts(pkgPath string, mode parser.Mode, filenames ...string) (*Package, error) {
	var fset = token.NewFileSet()

	files := make([]*File, len(filenames))
	for i, filename := range filenames {
		astFile, err := parser.ParseFile(fset, filename, nil, mode)
		if err != nil {
			return nil, tracer.Trace(err)
		}

		files[i] = NewFile(astFile)
	}

	return NewPackage(pkgPath, files), nil
}

// PackagesFromDirs returns list of packages.
// pkgDirs params is a map from package's path to the dir path.
func PackagesFromDirs(pkgDirs map[string]string) ([]*Package, error) {
	return PackagesFromDirsWithOpts(pkgDirs, nil, DefaultParseMode)
}

func PackagesFromDirsWithOpts(pkgDirs map[string]string, filter func(fs.FileInfo) bool, mode parser.Mode) ([]*Package, error) {
	packages := make([]*Package, len(pkgDirs))
	var i int
	for pkgPath, dir := range pkgDirs {
		pkg, err := SinglePackageFromDirWithOpts(pkgPath, dir, filter, mode)
		if err != nil {
			return nil, tracer.Trace(err)
		}
		packages[i] = pkg
		i++
	}

	return packages, nil
}

func NewFile(f *ast.File) *File {
	interfaces := make([]*Interface, 0)
	structs := make([]*Struct, 0)
	imports := make([]*Import, len(f.Imports))

	for i, imp := range f.Imports { // collect imports
		var name string
		if imp.Name != nil {
			name = imp.Name.Name // it could be the alias or _.
		}

		imports[i] = &Import{
			Name: name,
			Path: imp.Path.Value[1 : len(imp.Path.Value)-1], // remove '"' at the beginning and end of the path.
		}
	}

	for _, dec := range f.Decls { // collect structs and interfaces in the package.
		genDecl, ok := dec.(*ast.GenDecl)
		if !ok || len(genDecl.Specs) == 0 {
			continue
		}

		t, ok := genDecl.Specs[0].(*ast.TypeSpec)
		if !ok {
			continue
		}

		if _, ok := t.Type.(*ast.InterfaceType); ok {
			interfaces = append(interfaces, ifaceMeta(genDecl, t))
			continue
		}

		if _, ok := t.Type.(*ast.StructType); ok {
			structs = append(structs, structMeta(genDecl, t))
			continue
		}
	}
	return &File{
		PackageName: f.Name.Name,
		ImportMap:   importsMap(imports),
		Imports:     imports,
		Interfaces:  interfaces,
		Structs:     structs,
	}
}

func NewFileByName(filename string, mode parser.Mode) (*File, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filename, nil, mode)
	if err != nil {
		return nil, tracer.Trace(err)
	}
	return NewFile(f), nil
}
