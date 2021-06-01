package lg

import (
	"go/ast"
)

// fieldIsEmbeddedInterface returns true if the field is an embedded type declaration in an interface or a struct.
func fieldIsEmbeddedInterface(f *ast.Field) bool {
	ident, ok := f.Type.(*ast.Ident)
	if ok {
		_, ok := ident.Obj.Decl.(*ast.TypeSpec)
		return ok
	}
	return false
}

// formatNode returns node's type as string.
func formatNode(src []byte, node ast.Expr) string {
	return string(src[node.Pos()-1 : node.End()-1])
}

func IsError(Type string) bool {
	return Type == "error"
}

func mergeMaps(methods, other map[string]MethodMetadata) {
	if methods == nil {
		methods = make(map[string]MethodMetadata)
	}

	for k, v := range other {
		methods[k] = v
	}
}
