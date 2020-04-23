package hexa

import (
	"errors"
	"fmt"
	"github.com/Kamva/tracer"
)

type (
	// ID is the id of entities across the hexa packages.
	// This is because hexa does not want to be dependent
	// to specific type of id. (e.g mongo ObjectID, mysql integer,...)
	ID interface {
		fmt.Stringer

		// Validate specify provided id is valid or not.
		Validate(id interface{}) error

		// From convert provided value to its id.
		// From will returns error if provided value
		// can not convert to an native id.
		From(id interface{}) error

		// MustFrom Same as FromString but on occur error, it will panic.
		MustFrom(id interface{})

		// Val returns the native id value (e.g ObjectID in mongo, ...).
		Val() interface{}

		// IsEqual say that two hexa id are equal or not.
		IsEqual(ID) bool
	}

	// stringID implements the ID with string type (use as "guest user's id" or "service user id").
	stringID string
)

func (s *stringID) String() string {
	return string(*s)
}

func (s *stringID) Validate(id interface{}) error {
	sID := new(stringID)
	return sID.From(id)
}

// From function does not do anything for string ID type,
// implement to just satisfy the interface.
func (s *stringID) From(id interface{}) error {
	if id == nil {
		return tracer.Trace(errors.New("id value is nil"))
	}
	if strID, ok := id.(string); ok {
		*s = stringID(strID)
	}
	return tracer.Trace(errors.New("invalid ID type"))
}

// MustFrom function does not do anything for string ID type,
// implement to just satisfy the interface.
func (s *stringID) MustFrom(id interface{}) {
	if err := s.From(id); err != nil {
		panic(err)
	}
}

func (s *stringID) Val() interface{} {
	return string(*s)
}

func (s *stringID) IsEqual(id ID) bool {
	if id == nil || s == nil {
		return false
	}
	if val, ok := id.(*stringID); !ok || val.Val().(string) != s.Val().(string) {
		return false
	}
	return true
}

// NewStringID returns new hexa ID.
func NewStringID(id string) ID {
	strID := stringID(id)
	return &strID
}

// Assertion
var _ ID = new(stringID)
