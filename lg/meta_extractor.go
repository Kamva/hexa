package lg

import (
	"fmt"
	"go/ast"
	"reflect"
	"strings"
)

func ifaceMeta(decl *ast.GenDecl, t *ast.TypeSpec) *Interface {
	methods := make([]*Method, 0)
	var embedded []*EmbeddedField
	for _, m := range t.Type.(*ast.InterfaceType).Methods.List {
		//if its embedded interface, add it to the embedded list
		if isEmbeddedNode(m) {
			embedded = append(embedded, embeddedFieldMeta(m))
			continue
		}

		methods = append(methods, methodMeta(m))
	}

	return &Interface{
		Doc:         prepareComments(decl.Doc.Text()),
		Annotations: annotationsFromCommentGroup(decl.Doc),
		Name:        t.Name.Name,
		Embedded:    embedded,
		Methods:     methods,
	}
}

func structMeta(decl *ast.GenDecl, t *ast.TypeSpec) *Struct {
	var fields []*Field
	var embedded []*EmbeddedField

	for _, field := range t.Type.(*ast.StructType).Fields.List {
		//if its embedded interface, add it to the embedded list
		if isEmbeddedNode(field) {
			embedded = append(embedded, embeddedFieldMeta(field))
			continue
		}

		fields = append(fields, fieldMeta(field))
	}

	return &Struct{
		Doc:         prepareComments(decl.Doc.Text()),
		Annotations: annotationsFromCommentGroup(decl.Doc),
		Name:        t.Name.Name,
		Embedded:    embedded,
		Fields:      fields,
	}
}

func methodMeta(method *ast.Field) *Method {
	var params []*MethodParam
	var results []*MethodResult
	funcNode := method.Type.(*ast.FuncType)

	if funcNode.Params != nil {
		for _, param := range funcNode.Params.List {
			for _, paramName := range param.Names {
				p := MethodParam{
					Name: paramName.Name,
					Type: typeName(param.Type),
				}
				params = append(params, &p)
			}
		}
	}

	if funcNode.Results != nil {
		for _, result := range funcNode.Results.List {
			resultType := typeName(result.Type)

			// for unnamed result
			if len(result.Names) == 0 {
				r := MethodResult{
					Name: "",
					Type: resultType,
				}
				results = append(results, &r)
			}

			for _, resultName := range result.Names {
				r := MethodResult{
					Name: resultName.Name,
					Type: resultType,
				}
				results = append(results, &r)
			}
		}
	}

	return &Method{
		Doc:         prepareComments(method.Doc.Text()),
		Annotations: annotationsFromCommentGroup(method.Doc),
		Name:        method.Names[0].Name,
		Params:      params,
		Results:     results,
	}
}

func fieldMeta(field *ast.Field) *Field {
	var tag string
	if field.Tag != nil {
		tag = field.Tag.Value
	}

	return &Field{
		Doc:         prepareComments(field.Doc.Text()),
		Annotations: annotationsFromCommentGroup(field.Doc),
		Name:        field.Names[0].Name,
		Type:        typeName(field.Type),
		Tag:         reflect.StructTag(tag),
	}
}

func embeddedFieldMeta(field *ast.Field) *EmbeddedField {
	var tag string
	if field.Tag != nil {
		tag = field.Tag.Value
	}

	return &EmbeddedField{
		Doc:         prepareComments(field.Doc.Text()),
		Annotations: annotationsFromCommentGroup(field.Doc),
		Type:        typeName(field.Type),
		Tag:         reflect.StructTag(tag),
	}
}

// typeName returns the name of the type from a "Type"
// param of an ast Node.
// t could be *ast.Field.Type, ...
func typeName(t ast.Expr) string {
	switch v := t.(type) {
	case *ast.Ident:
		return v.Name
	case *ast.SelectorExpr:
		return fmt.Sprintf("%s.%s", v.X.(*ast.Ident).Name, v.Sel.Name)
	}

	return ""
}

// prepareComments prepares comments to use in templates as comments.
// You can use comments.Text() as input of this method on the *ast.CommentGroup type.
func prepareComments(comments string) string {
	return strings.Replace(strings.TrimSuffix(comments, "\n"), "\n", "\n// ", -1)
}
