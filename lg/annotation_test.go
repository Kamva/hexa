package lg

import (
	"go/ast"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
)

func Test_annotationsFromCommentGroup(t *testing.T) {
	g := []*ast.Comment{
		&ast.Comment{
			Slash: 0,
			Text:  "// regular comment",
		},
		&ast.Comment{
			Slash: 0,
			Text:  "//@abc `a:\"b\"`",
		},
	}

	annotations := annotationsFromCommentGroup(g)
	require.Equal(t, 1, len(annotations))
	a := annotations[0]
	assert.Equal(t, "abc", a.Name)
	assert.Equal(t, reflect.StructTag(`a:"b"`), a.Tag)
}

func TestAnnotations_Lookup(t *testing.T) {
	g := []*ast.Comment{
		&ast.Comment{
			Slash: 0,
			Text:  "//@abc `a:\"b\"`",
		},
	}

	annotations := annotationsFromCommentGroup(g)
	require.Equal(t, 1, len(annotations))
	assert.Nil(t, annotations.Lookup("123"))
	assert.Equal(t, &annotations[0], annotations.Lookup("abc"))

}
