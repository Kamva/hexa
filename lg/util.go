package lg

import (
	"go/ast"
	"io/ioutil"
	"os"

	"github.com/kamva/tracer"
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

func readAll(filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, tracer.Trace(err)
	}
	defer file.Close()

	src, err := ioutil.ReadAll(file)
	return src, tracer.Trace(err)
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
