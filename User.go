package hexa

import (
	"errors"
)

type (
	// Note: although getter function in Go Don't need to start with "Get" word, but because
	// most user models use this fields (Email,Phone,Name,...) as their database fields, we
	// add "Get" prefix to each getter method on this interface.
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

		// Username can be unique username,email,phone number or
		// everything else that can use as username.
		GetUsername() string

		// IsActive specify that user is active or no.
		IsActive() bool

		// PermissionsList returns the use permissions list to
		//use in RBAC access control services (like Gate).
		GetPermissionsList() []string
	}
	guestUser struct {
	}

	// user is default implementation of hexa User for real users.
	user struct {
		id       ID
		email    string
		phone    string
		name     string
		username string
		isActive bool
		perms    []string
	}

	// guestID is implementation of specific ID
	guestID string
)

// guestUserID is the guest user's id
var guestUserID = "__guest_id__"

func (u *user) IsGuest() bool {
	return false
}

func (u *user) Identifier() ID {
	return u.id
}

func (u *user) GetEmail() string {
	return u.email
}

func (u *user) GetPhone() string {
	return u.phone
}

func (u *user) GetName() string {
	return u.name
}

func (u *user) GetUsername() string {
	return u.email
}

func (u *user) IsActive() bool {
	return u.isActive
}

func (u *user) GetPermissionsList() []string {
	return u.perms
}

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

func (g guestUser) IsActive() bool {
	return false
}

func (g guestUser) GetPermissionsList() []string {
	return nil
}

// NewUser returns new hexa user instance.
func NewUser(id ID, email, phone, name, username string, isActive bool, perms []string) User {
	return &user{
		id:       id,
		email:    email,
		phone:    phone,
		name:     name,
		username: username,
		perms:    perms,
		isActive: isActive,
	}
}

func NewGuestUser() User {
	return guestUser{}
}

// Assert guestUser implements the User interface.
var _ ID = guestID("")
var _ User = &guestUser{}
var _ User = &user{}
