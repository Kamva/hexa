package hexa

import "fmt"

// ID is the id of entities across the hexa packages.
// This is because hexa does not want to be dependent
// to specific type of id. (e.g mongo ObjectID, mysql integer,...)
type ID interface {
	fmt.Stringer

	// Validate say that id value is valid or not.
	Validate(id interface{}) error

	// From convert provided value to its id.
	// From will returns error if provided value
	// can not convert to an native id.
	From(id interface{}) error

	// MustFrom Same as FromString but on occur error, it will panic.
	MustFrom(id interface{})

	// Val returns the native id value (e.g ObjectID in mongo, ...).
	Val() interface{}

	// Equal say that two hexa id are equal or not.
	Equal(ID) bool
}