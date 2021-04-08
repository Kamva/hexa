package lg

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"

	"github.com/kamva/tracer"
)

type InterfaceMetadata struct {
	Name    string
	Methods map[string]MethodMetadata
}

type MethodMetadata struct {
	Name    string
	Params  []MethodParam
	Results []MethodResult
}

type MethodParam struct {
	Name string
	Type string
}

type MethodResult struct {
	Name string
	Type string
}

func (r MethodResult) joinNameAndType() string {
	if r.Name == "" {
		return r.Type
	}
	return fmt.Sprintf("%s %s", r.Name, r.Type)
}

func ExtractInterfaceMetadata(srcfile string, ifaceName string) (*InterfaceMetadata, error) {
	fset := token.NewFileSet()
	src, err := readAll(srcfile)
	if err != nil {
		return nil, tracer.Trace(err)
	}

	f, err := parser.ParseFile(fset, srcfile, src, parser.AllErrors|parser.ParseComments)
	if err != nil {
		return nil, tracer.Trace(err)
	}

	metadata := InterfaceMetadata{
		Name:    ifaceName,
		Methods: make(map[string]MethodMetadata),
	}

	ast.Inspect(f, func(node ast.Node) bool {
		t, ok := node.(*ast.TypeSpec)
		if !ok || t.Name.Name != ifaceName {
			return true
		}

		for _, m := range t.Type.(*ast.InterfaceType).Methods.List {
			metadata.Methods[m.Names[0].Name] = extractMethodMetadata(src, m)
		}

		return false
	})

	return &metadata, nil
}

func extractMethodMetadata(src []byte, method *ast.Field) MethodMetadata {
	params := []MethodParam{}
	results := []MethodResult{}
	funcNode := method.Type.(*ast.FuncType)

	if funcNode.Params != nil {
		for _, param := range funcNode.Params.List {
			for _, paramName := range param.Names {
				p := MethodParam{
					Name: paramName.Name,
					Type: formatNode(src, param.Type),
				}
				params = append(params, p)
			}
		}
	}

	if funcNode.Results != nil {
		for _, result := range funcNode.Results.List {
			resultType := formatNode(src, result.Type)

			// for unnamed result
			if len(result.Names) == 0 {
				r := MethodResult{
					Name: "",
					Type: resultType,
				}
				results = append(results, r)
			}

			for _, resultName := range result.Names {
				r := MethodResult{
					Name: resultName.Name,
					Type: resultType,
				}
				results = append(results, r)
			}
		}
	}

	return MethodMetadata{
		Name:    method.Names[0].Name,
		Params:  params,
		Results: results,
	}
}
