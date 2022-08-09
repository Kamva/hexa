package lg

import (
	"fmt"
	"go/ast"
	"strings"

	"github.com/kamva/hexa/hlog"
)

// specialCharsBeforeType are the characters that may exist before a type.
// e.g., *string or []string or map[string]Health.
const specialCharsBeforeType = "]*"

// typeStr returns the type as string.
// t could be *ast.Field.Type, ...
func typeStr(t ast.Expr) string {
	switch v := t.(type) {
	case *ast.Ident:
		return v.Name
	case *ast.SelectorExpr:
		return fmt.Sprintf("%s.%s", v.X.(*ast.Ident).Name, v.Sel.Name)
	case *ast.StarExpr:
		return fmt.Sprintf("*%s", typeStr(v.X))
	case *ast.ArrayType:
		var length string
		if basicLit, ok := v.Len.(*ast.BasicLit); ok {
			length = basicLit.Value
		}
		return fmt.Sprintf("[%s]%s", length, typeStr(v.Elt))
	case *ast.MapType:
		return fmt.Sprintf("map[%s]%s", typeStr(v.Key), typeStr(v.Value))
	case *ast.ChanType:
		var send string
		var recv string

		if v.Dir == ast.SEND {
			send = "<-"
		}
		if v.Dir == ast.RECV {
			recv = "<-"
		}
		return fmt.Sprintf("%schan%s %s", recv, send, typeStr(v.Value))
	}

	hlog.Error("can not detect type to convert it to string", hlog.String("type", fmt.Sprintf("%T", t)))
	return ""
}

// parseType returns package name and the type.
// e.g., hexa.Health => returns "hexa","Health"
// e.g., Health => returns "","Health"
// e.g., *hexa.Health => returns "hexa","Health"
// Please note it just extracts the package name and the type name.
// so you can not detect if type is defined by pointer or is a map,...
// from this method results.
func parseType(t string) (string, string) {
	pkgStartIdx := strings.LastIndexAny(t, specialCharsBeforeType) + 1
	pkgEndIdx := strings.Index(t[pkgStartIdx:], ".")
	if pkgEndIdx == -1 { // If we don't have any dot, so we don't have any package, so just return the type.
		return "", t[pkgStartIdx:]
	}
	pkgEndIdx += pkgStartIdx

	return t[pkgStartIdx:pkgEndIdx], t[pkgEndIdx+1:]
}

// SetPackageOnType sets the package in the type.
// e.g., hexa, Health => hexa.Health
// e.g., hexa, []*Health => []*hexa.Health
func SetPackageOnType(pkg string, t string) string {
	if pkg != "" {
		pkg = pkg + "."
	}

	pkgStartIdx := strings.LastIndexAny(t, specialCharsBeforeType) + 1
	pkgEndIdx := strings.Index(t[pkgStartIdx:], ".") + pkgStartIdx

	return t[:pkgStartIdx] + pkg + t[pkgEndIdx+1:]
}

// isPrivateType returns true of the provided type is a private type.
func isPrivateType(t string) bool {
	_, t = parseType(t)
	return IsPrivate(t)
}
