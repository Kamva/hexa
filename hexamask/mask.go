package hexamask

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Kamva/gutil"
	"github.com/Kamva/hexa"
	"github.com/Kamva/tracer"
	"strings"
)

// FieldMask contains masked fields in input,output,...
// use it to detect which fields provided by user for in PATCH request.
type FieldMask struct {
	paths []string
}

func (m *FieldMask) UnmarshalJSON(b []byte) error {
	// Check null values
	if string(b) == "null" {
		return nil
	}
	// Check invalid values
	if len(b) < 2 || b[0] != '"' || b[len(b)-1] != '"' {
		return tracer.Trace(errors.New("hexa mask: bad JSON key"))
	}

	val := string(b[1 : len(b)-1])
	m.SetPaths([]string{})
	if len(val) != 0 {
		m.SetPaths(strings.Split(val, ","))
	}
	return nil
}

func (m *FieldMask) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, strings.Join(m.paths, ","))), nil
}

func (m *FieldMask) SetPaths(paths []string) {
	m.paths = paths
}

// IsMasked specifies whether the provided path is masked or not.
func (m *FieldMask) IsMasked(path string) bool {
	return gutil.Contains(m.paths, path)
}

func (m *FieldMask) String() string {
	return strings.Join(m.paths, ",")
}

// MaskMakPaths mask all paths in the provided map with the provided depth.
func MaskMakPaths(m hexa.Map, mask *FieldMask, depth int) {
	mask.SetPaths(gutil.MapPathExtractor{Depth: depth, Separator: "."}.Extract(m))
}

var _ fmt.Stringer = &FieldMask{}
var _ json.Unmarshaler = &FieldMask{}
var _ json.Marshaler = &FieldMask{}
