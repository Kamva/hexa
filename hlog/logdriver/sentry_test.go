package logdriver

import (
	"errors"
	"testing"

	"github.com/kamva/hexa/hlog"
)

func TestExtractError(t *testing.T) {
	wrapped := errors.New("db connection refused")

	tests := []struct {
		name string
		args []hlog.Field
		want error
	}{
		{
			name: "no fields",
			args: nil,
			want: nil,
		},
		{
			name: "no error field",
			args: []hlog.Field{hlog.String("user_id", "42"), hlog.Int("count", 3)},
			want: nil,
		},
		{
			name: "error field via hlog.Err",
			args: []hlog.Field{hlog.Any("user_id", "42"), hlog.Err(wrapped)},
			want: wrapped,
		},
		{
			name: "error field via hlog.NamedErr",
			args: []hlog.Field{hlog.NamedErr("cause", wrapped)},
			want: wrapped,
		},
		{
			name: "nil error field",
			args: []hlog.Field{hlog.Err(nil)},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := extractError(tt.args); got != tt.want {
				t.Fatalf("extractError() = %v, want %v", got, tt.want)
			}
		})
	}
}
