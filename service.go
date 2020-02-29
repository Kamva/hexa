package kitty

// LayeredService is specific type of service that is layered relative to functionality.
// e.g validation->sanitizer->guard->core (each service has next service (of same type) and a root (first service))
type LayeredService interface {
	SetRoot(root interface{})
	SetNext(next interface{})
}

// SetServiceChain set services chain in order of provided services.
// finally returns the root of services.
func SetServiceChain(services ...LayeredService) interface{} {
	if len(services) == 0 {
		return nil
	}
	root := services[0]

	var prev LayeredService
	for _, current := range services {
		current.SetRoot(root)

		if prev != nil {
			prev.SetNext(current)
		}

		prev = current
	}

	return root
}
