package hlog

import (
	"go.uber.org/zap/zapcore"
)

func fieldToKeyVal(f Field) (key string, val any) {
	switch f.Type {
	case zapcore.Int64Type, zapcore.Int32Type, zapcore.Int16Type, zapcore.Int8Type, zapcore.UintptrType,
		zapcore.Uint64Type, zapcore.Uint32Type, zapcore.Uint16Type, zapcore.Uint8Type:
		val = f.Integer
	case zapcore.StringType:
		val = f.String
	default:
		val = f.Interface
	}

	return f.Key, val
}

func fieldsToMap(fields ...Field) map[string]any {
	m := make(map[string]any)
	for _, f := range fields {
		k, v := fieldToKeyVal(f)
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
