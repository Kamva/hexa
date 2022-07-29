package hexamask

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/kamva/gutil"
	"github.com/kamva/hexa"
	"github.com/kamva/tracer"
)

// FieldMask mask keep list of masked paths of a struct|map|...
// and then you can check whether a path|struct field is masked or not.
// Use it to detect which fields provided by a user in a PATCH request.
// Note: you can specify mask path of each field of a struct by set the
// "mask" tag, otherwise it check the "json" tag.
type FieldMask struct {
	paths        []string
	maskedFields []any
}

func (fm *FieldMask) UnmarshalJSON(b []byte) error {
	// Check null values
	if string(b) == "null" {
		return nil
	}
	// Check invalid values
	if len(b) < 2 || b[0] != '"' || b[len(b)-1] != '"' {
		return tracer.Trace(errors.New("hexa mask: bad JSON key"))
	}

	val := string(b[1 : len(b)-1])
	fm.SetPaths([]string{})
	if len(val) != 0 {
		fm.SetPaths(strings.Split(val, ","))
	}
	return nil
}

func (fm *FieldMask) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, strings.Join(fm.paths, ","))), nil
}

// SetPaths sets the masked paths
func (fm *FieldMask) SetPaths(paths []string) {
	fm.paths = paths
}

// PathIsMasked tel you whether the provided path is masked or not.
func (fm *FieldMask) PathIsMasked(path string) bool {
	return gutil.Contains(fm.paths, path)
}

// IsMasked gets a struct field and specifies whether that
// field is masked or not. before call to this method you
// must call to the Mask() method to mask a struct.
// Note: provided value must be pointer to the value.
// even if the field is a pointer, you must provide
// pointer to that pointer field.
func (fm *FieldMask) IsMasked(i any) bool {
	if reflect.TypeOf(i).Kind() != reflect.Ptr {
		panic("value must be interface")
	}
	if fm.maskedFields == nil {
		panic("you must mask a struct before invoking the \"IsMasked\" method")
	}
	for _, f := range fm.maskedFields {
		if f == i {
			return true
		}
	}
	return false
}

// Mask masks a struct's fields and then you can check
// to detect whether a field of that struct is masked or not.
// Note: provided value must be pointer to a struct.
func (fm *FieldMask) Mask(s any) {
	v := reflect.ValueOf(s)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		panic("provided value must be a pointer to a struct")
	}

	fm.maskedFields = fm.maskStruct("", v)
}

// maskStruct get a reflect Value of pointer to struct and
// add all masked fields of the struct to the masked fields
// list.
// Note: provided value must be a pointer of a value otherwise
// it panic.
func (fm *FieldMask) maskStruct(pathPrefix string, v reflect.Value) []any {
	if v.IsNil() || v.Elem().Kind() != reflect.Struct {
		return nil
	}
	iv := v.Elem()
	it := iv.Type()
	maskList := make([]any, 0)
	for i := 0; i < it.NumField(); i++ {
		fieldValue := iv.Field(i)

		path := fm.pathOf(it.Field(i).Tag)
		// if mask value of a field is not specified, we must ignore it
		if path == "" {
			continue
		}
		if pathPrefix != "" {
			path = fmt.Sprintf("%s.%s", pathPrefix, path)
		}
		if fm.PathIsMasked(path) {
			maskList = append(maskList, fieldValue.Addr().Interface())
		}

		maskList = append(maskList, fm.maskStruct(path, fieldValue.Addr())...)
	}

	return maskList
}

// String returns joined paths which divided by "," rune.
func (fm *FieldMask) String() string {
	return strings.Join(fm.paths, ",")
}

// pathOf gets a StructTag and extracts the mask path
// from the "mask" or "json" tag values.
func (fm *FieldMask) pathOf(tag reflect.StructTag) string {
	if v, ok := tag.Lookup("mask"); ok {
		return v
	}
	return tag.Get("json")
}

// MaskMapPaths mask all paths in the provided map with the provided depth.
func MaskMapPaths(m hexa.Map, mask *FieldMask, depth int) {
	mask.SetPaths(gutil.MapPathExtractor{Depth: depth, Separator: "."}.Extract(m))
}

var _ fmt.Stringer = &FieldMask{}
var _ json.Unmarshaler = &FieldMask{}
var _ json.Marshaler = &FieldMask{}
