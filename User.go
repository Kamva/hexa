package hexa

import (
	"errors"
)

type (
	User interface {
		// Specify that user is guestUser or no.
		IsGuest() bool

		// Return users identifier (if guestUser return just empty string or something like this.)
		Identifier() ID

		// GetEmail returns the user's email.
		// This value can be empty.
		GetEmail() string

		// GetPhone returns the user's phone number.
		// This value can be empty.
		GetPhone() string

		// Return the user name.
		GetName() string

		// Username can be unique username,email,phone number or everything else that can use as username.
		GetUsername() string
	}
	guestUser struct {
	}

	// guestID is implementation of specific ID
	guestID string
)

// guestUserID is the guest user's id
var guestUserID = "__guest_id__"

func (g guestID) Validate(id interface{}) error {
	if idStr, ok := id.(string); ok && idStr == guestUserID {
		return nil
	}

	return errors.New("guest user id is not valid")
}

func (g guestID) String() string {
	return string(g)
}

// From function does not do anything for guest ID type,
// implement to just satisfy the interface.
func (g guestID) From(id interface{}) error {
	return nil
}

// MustFrom function does not do anything for guest ID type,
// implement to just satisfy the interface.
func (g guestID) MustFrom(id interface{}) {
	// empty
}

func (g guestID) Equal(hexaID ID) bool {
	if hexaID == nil {
		return false
	}

	_, ok := hexaID.Val().(guestID)

	return ok
}

func (g guestID) Val() interface{} {
	return string(g)
}

func (g guestUser) IsGuest() bool {
	return true
}

func (g guestUser) Identifier() ID {
	return guestID(guestUserID)
}

func (g guestUser) GetEmail() string {
	return ""
}

func (g guestUser) GetPhone() string {
	return ""
}

func (g guestUser) GetName() string {
	return "__guest__"
}

func (g guestUser) GetUsername() string {
	return "__guest__username__"
}

func NewGuestUser() User {
	return guestUser{}
}

// Assert guestUser implements the User interface.
var _ ID = guestID("")
var _ User = &guestUser{}
