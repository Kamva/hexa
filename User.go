package kitty

type User interface {
	// Specify that user is guestUser or no.
	IsGuest() bool

	// Return users identifier (if guestUser return just empty string or something like this.)
	GetID() interface{}

	// Username can be unique username,email,phone number or everything else that can use as username.
	GetUsername() string
}

type guestUser struct {
}

func (g guestUser) IsGuest() bool {
	return true
}

func (g guestUser) GetID() interface{} {
	return "__guest_id__"
}

func (g guestUser) GetUsername() string {
	return "__guest__username__"
}

func NewGuestUser() User {
	return guestUser{}
}

// Assert guestUser implements the User interface.
var _ User = &guestUser{}
