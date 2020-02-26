package kitty

import "encoding/json"

type Secret string

func (s Secret) String() string {
	return "****"
}

func (s Secret) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

