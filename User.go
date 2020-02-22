package kitty

type User interface {
	// Specify that user is guest or nop.
	IsGuest() bool

	// Return users identifier (if guest return just empty string or something like this.)
	GetID() interface{}

	// Username can be unique username,email,phone number or everything else that can use as username.
	GetUsername() string
}
