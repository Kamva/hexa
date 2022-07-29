package hexa

import "encoding/json"

// Use secret to string show as * in fmt package.
type Secret string

// Map defines a well-known Golang map: map[string]any
type Map map[string]any

func (s Secret) String() string {
	return "****"
}

func (s Secret) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}
