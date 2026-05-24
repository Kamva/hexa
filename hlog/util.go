package hlog

import (
	"go.uber.org/zap/zapcore"
)

func FieldToKeyVal(f Field) (key string, val any) {
	// Let zap decode the field into its real Go value. The previous manual
	// logic returned f.Integer for any non-string/non-interface field, which
	// is wrong for bool/float/duration/time fields (zap packs those into the
	// Integer bits), and returned the Stringer object instead of its string.
	enc := zapcore.NewMapObjectEncoder()
	f.AddTo(enc)
	return f.Key, enc.Fields[f.Key]
}

func fieldsToMap(fields ...Field) map[string]any {
	m := make(map[string]any)
	for _, f := range fields {
		k, v := FieldToKeyVal(f)
		m[k] = v
	}
	return m
}

func MapToFields(m map[string]any) []Field {
	fields := make([]Field, 0)
	for k, v := range m {
		fields = append(fields, Any(k, v))
	}
	return fields
}
