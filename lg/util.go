package lg

import (
	"go/ast"
	"path"
	"reflect"
	"runtime"
	"strings"
	"unicode"
	"unicode/utf8"
)

// IsPrivate returns ture if the value is a private value.
// As you know, a type, field, interface,... is private in
// go if it's first letter is a lowercase.
func IsPrivate(val string) bool {
	r, _ := utf8.DecodeRuneInString(val) // get the first rune not the first byte.
	return r != utf8.RuneError && unicode.IsLower(r)
}

// isEmbeddedNode returns true if the field is an embedded
// type declaration in an interface or a struct.
func isEmbeddedNode(f *ast.Field) bool {
	// If a field doesn't have any name, so it's an
	// embedded field.
	return len(f.Names) == 0
}

func IsError(Type string) bool {
	return Type == "error"
}

// importsMap maps the package's name or alias to the package path.
func importsMap(l []*Import) map[string]string {
	m := make(map[string]string)
	for _, imp := range l {
		if imp.Name == "_" { // ignore blank imports.
			continue
		}

		if imp.Name == "" {
			m[path.Base(imp.Path)] = imp.Path
			continue
		}
		m[imp.Name] = imp.Path // for aliases
	}
	return m
}

func IsSamePackage(pkg *Package, target *Package) bool {
	return pkg.Path == target.Path
}

func FnName(fn any) string {
	name := runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
	return name[strings.LastIndex(name, ".")+1:]
}

func Lookup(tag reflect.StructTag, keys ...string) (string, bool) {
	for _, k := range keys {
		if v, ok := tag.Lookup(k); ok {
			return v, ok
		}
	}
	return "", false
}
