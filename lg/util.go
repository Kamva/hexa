package lg

import (
	"fmt"
	"go/ast"
	"path"
	"strings"
)

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

// parseType returns package name and the type.
// e.g., hexa.Health => returns "hexa","Health"
// e.g., Health => returns "","Health"
func parseType(t string) (string, string) {
	idx := strings.Index(t, ".")
	if idx == -1 {
		return "", t
	}

	return t[0:idx], t[idx+1:]
}

// constructType creates type from its package name and the type name.
// e.g., hexa, Health =>
func constructType(pkg string, t string) string {
	if pkg == "" {
		return t
	}

	return fmt.Sprintf("%s.%s", pkg, t)
}
