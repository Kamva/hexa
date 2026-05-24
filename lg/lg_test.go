package lg

import (
	"go/parser"
	"go/token"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsPrivateAndError(t *testing.T) {
	assert.True(t, IsPrivate("health"))
	assert.False(t, IsPrivate("Health"))
	assert.False(t, IsPrivate(""))

	assert.True(t, IsError("error"))
	assert.False(t, IsError("Err"))
}

func TestUseTypeInPackage(t *testing.T) {
	from := &Package{Name: "hexa"}
	assert.Equal(t, "hexa.Health", UseTypeInPackage(from, "Health")) // local exported type gets qualified
	assert.Equal(t, "string", UseTypeInPackage(from, "string"))      // primitives are left alone
	assert.Equal(t, "other.X", UseTypeInPackage(from, "other.X"))    // already-qualified is left alone
}

func TestLookup(t *testing.T) {
	tag := reflect.StructTag(`json:"id" mask:"identifier"`)
	v, ok := Lookup(tag, "mask", "json")
	assert.True(t, ok)
	assert.Equal(t, "identifier", v)

	_, ok = Lookup(tag, "missing")
	assert.False(t, ok)
}

func TestRunLayerFns(t *testing.T) {
	var order []string
	fns := map[string]func() error{
		"a": func() error { order = append(order, "a"); return nil },
		"b": func() error { order = append(order, "b"); return nil },
	}
	require.NoError(t, RunLayerFns(fns, "b", "a"))
	assert.Equal(t, []string{"b", "a"}, order)
}

func TestTypeStr(t *testing.T) {
	cases := map[string]string{
		"*Foo":           "*Foo",
		"[]Foo":          "[]Foo",
		"[3]Foo":         "[3]Foo",
		"map[string]Foo": "map[string]Foo",
		"chan int":       "chan int",
		"<-chan int":     "<-chan int",
		"pkg.Foo":        "pkg.Foo",
	}
	for in, want := range cases {
		expr, err := parser.ParseExpr(in)
		require.NoError(t, err, in)
		assert.Equal(t, want, typeStr(expr), in)
	}
}

func TestNewFile_ExtractsMetadata(t *testing.T) {
	src := "package sample\n" +
		"import (\n" +
		"\t\"context\"\n" +
		"\thx \"github.com/kamva/hexa\"\n" +
		")\n" +
		"// Doer does things.\n" +
		"type Doer interface {\n" +
		"\thx.Health\n" +
		"\tDo(ctx context.Context, n int) (string, error)\n" +
		"}\n" +
		"type Model struct {\n" +
		"\tID   string `json:\"id\"`\n" +
		"\tName string\n" +
		"}\n"

	fset := token.NewFileSet()
	astFile, err := parser.ParseFile(fset, "sample.go", src, parser.ParseComments)
	require.NoError(t, err)

	f := NewFile(astFile)
	assert.Equal(t, "sample", f.PackageName)
	assert.Equal(t, "context", f.ImportMap["context"])
	assert.Equal(t, "github.com/kamva/hexa", f.ImportMap["hx"])

	iface := f.FindInterface("Doer")
	require.NotNil(t, iface)
	require.Len(t, iface.Embedded, 1)
	assert.Equal(t, "hx.Health", iface.Embedded[0].Type)

	do := iface.MethodByName("Do")
	require.NotNil(t, do)
	require.Len(t, do.Params, 2)
	assert.Equal(t, "context.Context", do.Params[0].Type)
	assert.Equal(t, "int", do.Params[1].Type)
	require.Len(t, do.Results, 2)
	assert.Equal(t, "string", do.Results[0].Type)
	assert.Equal(t, "error", do.Results[1].Type)

	model := f.FindStruct("Model")
	require.NotNil(t, model)
	require.Len(t, model.Fields, 2)
	assert.Equal(t, "ID", model.Fields[0].Name)
	assert.Equal(t, "id", model.Fields[0].Tag.Get("json"))
}

func TestEmbeddedResolver_MergesEmbeddedInterfaceMethods(t *testing.T) {
	src := "package p\n" +
		"type A interface { Foo() error }\n" +
		"type B interface {\n\tA\n\tBar()\n}\n"

	fset := token.NewFileSet()
	astFile, err := parser.ParseFile(fset, "p.go", src, parser.ParseComments)
	require.NoError(t, err)

	pkg := NewPackage("example.com/p", []*File{NewFile(astFile)})
	require.NoError(t, NewEmbeddedResolver(pkg).Resolve())

	_, b := pkg.FindInterface("B")
	require.NotNil(t, b)
	assert.NotNil(t, b.MethodByName("Bar"))
	assert.NotNil(t, b.MethodByName("Foo")) // merged from embedded A
	assert.Len(t, b.Methods, 2)
}

func TestUseEmbeddedFieldsInPackage_ReturnsRepackaged(t *testing.T) {
	from := &Package{Name: "hexa"}
	out := UseEmbeddedFieldsInPackage(from, []*EmbeddedField{{Type: "Health"}})
	require.Len(t, out, 1)
	// Regression: used to return the original (unqualified) slice.
	assert.Equal(t, "hexa.Health", out[0].Type)
}

func TestFuncs_ResultsHelpers(t *testing.T) {
	funcs := Funcs()

	hasErr := funcs["hasErrInResults"].(func([]*MethodResult) bool)
	errVar := funcs["errResultVar"].(func([]*MethodResult) string)
	withErr := []*MethodResult{{Type: "string"}, {Type: "error"}}
	withoutErr := []*MethodResult{{Type: "string"}}

	assert.True(t, hasErr(withErr))
	assert.False(t, hasErr(withoutErr))
	assert.Equal(t, "r2", errVar(withErr))
	assert.Equal(t, "r1", ResultVar(0))

	title := funcs["title"].(func(string) string)
	assert.Equal(t, "Foo", title("foo"))
}
