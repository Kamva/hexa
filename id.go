package kitty

import "fmt"

// ID is the id of entities across the kitty packages.
// This is because kitty does not want to be dependent
// to specific type of id. (e.g mongo ObjectID, mysql integer,...)
type ID interface {
	fmt.Stringer

	// From convert provided value to its id.
	// From will returns error if provided value
	// can not convert to an native id.
	From(id interface{}) error

	// MustFrom Same as FromString but on occur error, it will panic.
	MustFrom(id interface{})

	// Val returns the native id value (e.g ObjectID in mongo, ...).
	Val() interface{}
}
