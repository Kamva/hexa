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
		{
			Slash: 0,
			Text:  "// regular comment",
		},
		{
			Slash: 0,
			Text:  "// @abc",
		},
		{
			Slash: 0,
			Text:  "//@cde `a:\"b\"`",
		},
	}

	annotations := annotationsFromCommentGroup(&ast.CommentGroup{List: g})
	require.Equal(t, 2, len(annotations))

	first := annotations[0]
	assert.Equal(t, "abc", first.Name)

	second := annotations[1]
	assert.Equal(t, "cde", second.Name)
	assert.Equal(t, reflect.StructTag(`a:"b"`), second.Tag)
}

func TestAnnotations_Lookup(t *testing.T) {
	g := []*ast.Comment{
		{
			Slash: 0,
			Text:  "//@abc `a:\"b\"`",
		},
	}

	annotations := annotationsFromCommentGroup(&ast.CommentGroup{List: g})
	require.Equal(t, 1, len(annotations))
	assert.Nil(t, annotations.Lookup("123"))
	assert.Equal(t, &annotations[0], annotations.Lookup("abc"))

}
