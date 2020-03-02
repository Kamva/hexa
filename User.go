package kitty

type (
	User interface {
		// Specify that user is guestUser or no.
		IsGuest() bool

		// Return users identifier (if guestUser return just empty string or something like this.)
		Identifier() ID

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

func (g guestID) IsValid(id interface{}) bool {
	if idStr, ok := id.(string); ok {
		return idStr == guestUserID
	}

	return false
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

func (g guestID) Val() interface{} {
	return string(g)
}

func (g guestUser) IsGuest() bool {
	return true
}

func (g guestUser) Identifier() ID {
	return guestID(guestUserID)
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
