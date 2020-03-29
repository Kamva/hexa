package hexa

import (
	"errors"
	"github.com/Kamva/gutil"
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

	// user is default implementation of hexa User for real users.
	user struct {
		id       ID
		isGuest  bool
		email    string
		phone    string
		name     string
		username string
		isActive bool
		perms    []string
	}

	// guestID is implementation of specific ID
	guestID string

	// UserExporter export a user to json and then import it.
	ExportedUser struct {
		ID       interface{} `json:"id"`
		IsGuest  bool        `json:"is_guest"`
		Email    string      `json:"email"`
		Phone    string      `json:"phone"`
		Name     string      `json:"name"`
		Username string      `json:"username"`
		IsActive bool        `json:"is_active"`
		Perms    []string    `json:"perms"`
	}

	// UserExporter export a user to json and then import it.
	UserExporter struct{}
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

// Export method export a user to map.
func (e UserExporter) Export(user User) map[string]interface{} {
	return gutil.StructToMap(ExportedUser{
		ID:       user.Identifier().Val(),
		IsGuest:  user.IsGuest(),
		Email:    user.GetEmail(),
		Phone:    user.GetPhone(),
		Name:     user.GetName(),
		Username: user.GetUsername(),
		IsActive: user.IsActive(),
		Perms:    user.GetPermissionsList(),
	})
}

// Import method a user from map.
func (e UserExporter) Import(id ID, exportedMap map[string]interface{}) (User, error) {
	eu := ExportedUser{}
	err := gutil.MapToStruct(exportedMap, &eu)
	if err != nil {
		return nil, err
	}

	if eu.IsGuest {
		id = guestID(guestUserID)
	} else {
		if err := id.From(eu.ID); err != nil {
			return nil, err
		}
	}

	user := NewUser(id, eu.Email, eu.Phone, eu.Name, eu.Username, eu.IsActive, eu.Perms)

	return user, nil
}

// NewUser returns new hexa user instance.
func NewUser(id ID, email, phone, name, username string, isActive bool, perms []string) User {
	return &user{
		id:       id,
		isGuest:  id.Equal(guestID(guestUserID)),
		email:    email,
		phone:    phone,
		name:     name,
		username: username,
		isActive: isActive,
		perms:    perms,
	}
}

// NewGuestUser returns new instance of the guest user.
func NewGuestUser() User {
	email := ""
	phone := ""
	name := "__guest__"
	username := "__guest__username__"
	return NewUser(guestID(guestUserID), email, phone, name, username, false, []string{})
}

// Assert guestUser implements the User interface.
var _ ID = guestID("")
var _ User = &user{}
