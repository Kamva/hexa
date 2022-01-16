package lg

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"

	"github.com/kamva/gutil"
	"github.com/kamva/tracer"
)

type InterfaceMetadata struct {
	Doc     string
	Doc2    string
	Name    string
	Methods map[string]MethodMetadata
}

type MethodMetadata struct {
	Doc         string
	Annotations Annotations
	Name        string
	Params      []MethodParam
	Results     []MethodResult
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

// ExtractInterfaceMetadata extracts the interface meta data, it also recursively extract
// embedded interfaces metadata.
func ExtractInterfaceMetadata(srcfile string, ifaceName string) (*InterfaceMetadata, error) {
	fset := token.NewFileSet()
	src, err := gutil.ReadAll(srcfile)
	if err != nil {
		return nil, tracer.Trace(err)
	}

	f, err := parser.ParseFile(fset, srcfile, src, parser.AllErrors|parser.ParseComments)
	if err != nil {
		return nil, tracer.Trace(err)
	}

	var doc string
	if cg := extractDocFromType(f, ifaceName); cg != nil {
		doc = prepareComments(cg.Text())
	}

	return &InterfaceMetadata{
		Name:    ifaceName,
		Doc:     doc,
		Methods: extractInterfaceMethods(src, f, ifaceName),
	}, nil
}

func extractDocFromType(f *ast.File, typeName string) *ast.CommentGroup {
	var doc *ast.CommentGroup
	ast.Inspect(f, func(node ast.Node) bool {
		if decl, ok := node.(*ast.GenDecl); ok && len(decl.Specs) != 0 {
			if t, ok := decl.Specs[0].(*ast.TypeSpec); ok && t.Name.Name == typeName {
				doc = decl.Doc
				return false
			}
		}

		return true
	})

	return doc
}

func extractInterfaceMethods(src []byte, f *ast.File, ifaceName string) map[string]MethodMetadata {
	methods := make(map[string]MethodMetadata)
	ast.Inspect(f, func(node ast.Node) bool {
		t, ok := node.(*ast.TypeSpec)
		if !ok || t.Name.Name != ifaceName {
			return true
		}

		for _, m := range t.Type.(*ast.InterfaceType).Methods.List {
			//if its embedded interface in parent, we need to extract its methods too.
			if fieldIsEmbeddedInterface(m) {
				mergeMaps(methods, extractInterfaceMethods(src, f, m.Type.(*ast.Ident).Name))
				continue
			}

			methods[m.Names[0].Name] = extractMethodMetadata(src, m)
		}

		return false
	})

	return methods
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
	var annotations Annotations
	if method.Doc != nil {
		annotations = annotationsFromCommentGroup(method.Doc.List)
	}
	return MethodMetadata{
		Doc:         prepareComments(method.Doc.Text()),
		Annotations: annotations,
		Name:        method.Names[0].Name,
		Params:      params,
		Results:     results,
	}
}

// prepareComments prepares comments to use in templates as comments.
// You can use comments.Text() as input of thie smethod on the *ast.CommentGroup type.
func prepareComments(comments string) string {
	return strings.Replace(strings.TrimSuffix(comments, "\n"), "\n", "\n// ", -1)
}
