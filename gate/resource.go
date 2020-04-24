package gate

import "github.com/Kamva/hexa"

// emptyResource satisfies the ResourceWithOwner interface, you can use
// it to check current user's id is equal to user_id payload or not.
// e.g If you want to create new post for somebody, you get the user_id
// as payload, so you need to check current user's id is equal to provided
// user_id in the payload or not.
type emptyResource struct {
	userID hexa.ID
}

func (e emptyResource) GateCheckOwnerIS(id hexa.ID) bool {
	return e.userID.IsEqual(id)
}

// NewEmptyResourceWithOwner returns new instance of the ResourceWithOwner
func NewEmptyResourceWithOwner(userID hexa.ID) ResourceWithOwner {
	return emptyResource{userID: userID}
}

// Assertion
var _ ResourceWithOwner = emptyResource{}
