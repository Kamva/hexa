package kitty

import "encoding/json"

// Use secret to string show as * in fmt package.
type Secret string

type Map map[string]interface{}

func (s Secret) String() string {
	return "****"
}

func (s Secret) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}


