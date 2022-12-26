package hexa

import "encoding/json"

// Map defines a well-known Golang map: map[string]any
type Map = map[string]any

// Secret is used to show * instead of the real string's
// content the in the fmt package.
type Secret string

func (s Secret) String() string {
	return "****"
}

func (s Secret) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}
