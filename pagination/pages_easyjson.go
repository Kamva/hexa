// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package pagination

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjsonD6eec276DecodeGithubComKamvaHexaPagination(in *jlexer.Lexer, out *Pagination) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "page":
			out.Page = int(in.Int())
		case "per_page":
			out.PerPage = int(in.Int())
		case "page_count":
			out.PageCount = int(in.Int())
		case "total_count":
			out.TotalCount = int(in.Int())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonD6eec276EncodeGithubComKamvaHexaPagination(out *jwriter.Writer, in Pagination) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"page\":"
		out.RawString(prefix[1:])
		out.Int(int(in.Page))
	}
	{
		const prefix string = ",\"per_page\":"
		out.RawString(prefix)
		out.Int(int(in.PerPage))
	}
	{
		const prefix string = ",\"page_count\":"
		out.RawString(prefix)
		out.Int(int(in.PageCount))
	}
	{
		const prefix string = ",\"total_count\":"
		out.RawString(prefix)
		out.Int(int(in.TotalCount))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Pagination) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD6eec276EncodeGithubComKamvaHexaPagination(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Pagination) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD6eec276EncodeGithubComKamvaHexaPagination(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Pagination) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD6eec276DecodeGithubComKamvaHexaPagination(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Pagination) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD6eec276DecodeGithubComKamvaHexaPagination(l, v)
}
func easyjsonD6eec276DecodeGithubComKamvaHexaPagination1(in *jlexer.Lexer, out *Pages) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "pagination":
			(out.Pagination).UnmarshalEasyJSON(in)
		case "items":
			if m, ok := out.Items.(easyjson.Unmarshaler); ok {
				m.UnmarshalEasyJSON(in)
			} else if m, ok := out.Items.(json.Unmarshaler); ok {
				_ = m.UnmarshalJSON(in.Raw())
			} else {
				out.Items = in.Interface()
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonD6eec276EncodeGithubComKamvaHexaPagination1(out *jwriter.Writer, in Pages) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"pagination\":"
		out.RawString(prefix[1:])
		(in.Pagination).MarshalEasyJSON(out)
	}
	{
		const prefix string = ",\"items\":"
		out.RawString(prefix)
		if m, ok := in.Items.(easyjson.Marshaler); ok {
			m.MarshalEasyJSON(out)
		} else if m, ok := in.Items.(json.Marshaler); ok {
			out.Raw(m.MarshalJSON())
		} else {
			out.Raw(json.Marshal(in.Items))
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Pages) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD6eec276EncodeGithubComKamvaHexaPagination1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Pages) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD6eec276EncodeGithubComKamvaHexaPagination1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Pages) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD6eec276DecodeGithubComKamvaHexaPagination1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Pages) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD6eec276DecodeGithubComKamvaHexaPagination1(l, v)
}
