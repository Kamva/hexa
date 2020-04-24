package gate

import (
	"errors"
	"github.com/Kamva/gutil"
	"github.com/Kamva/hexa"
)

// ResourceWithOwner specify user is owner of specific resource or not.
type ResourceWithOwner interface {
	// GateCheckOwnerIs specifies provided id is id of the resource owner or not.
	GateCheckOwnerIS(hexa.ID) bool
}

// UserOwnsResourcePolicy policy returns true if the user own provided resource.
func UserOwnsResourcePolicy(c hexa.Context, resource interface{}) (bool, error) {
	if gutil.IsNil(resource) {
		return false, nil
	}
	u := c.User()
	if m, ok := resource.(ResourceWithOwner); ok {
		return m.GateCheckOwnerIS(u.Identifier()), nil
	}
	return false, errors.New("provided resource is invalid")
}

// TruePolicy always returns true
func TruePolicy(c hexa.Context, r interface{}) (bool, error) {
	return true, nil
}

// DefaultPolicy is default policy for gates.
var DefaultPolicy = UserOwnsResourcePolicy

// Assertion
var _ hexa.GatePolicy = UserOwnsResourcePolicy
var _ hexa.GatePolicy = TruePolicy
