package kitty

// LayeredService is specific type of service that is layered relative to functionality.
// e.g validation->sanitizer->guard->core (each service has next service (of same type) and a root (first service))
type LayeredService interface {
	SetRoot(root interface{})
	SetNext(next interface{})
}
