package gate

import (
	"errors"
	"github.com/Kamva/gutil"
	"github.com/Kamva/hexa"
)

// UserWithOwner define model that has owner.
type ResourceWithOwner interface {
	// GateCheckOwnerIs method returns id.
	// content of this method can be something like
	// resourceHexaID.Equal(id)
	GateCheckOwnerIS(hexa.ID) bool
}

// UserOwnsResourcePolicy policy returns true if the user own provided resource.
func UserOwnsResourcePolicy(c hexa.Context, u hexa.User, r interface{}) (bool, error) {
	if gutil.IsNil(r) {
		return false, nil
	}

	if m, ok := r.(ResourceWithOwner); ok {
		return m.GateCheckOwnerIS(u.Identifier()), nil
	}
	return false, errors.New("provided resource is invalid")
}

// TruePolicy always returns true
func TruePolicy(c hexa.Context, u hexa.User, r interface{}) (bool, error) {
	return true, nil
}

// DefaultPolicy is default policy for gates.
var DefaultPolicy = UserOwnsResourcePolicy

// Assertion
var _ hexa.GatePolicy = UserOwnsResourcePolicy
