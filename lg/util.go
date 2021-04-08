package lg

import (
	"go/ast"
	"io/ioutil"
	"os"

	"github.com/kamva/tracer"
)

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
