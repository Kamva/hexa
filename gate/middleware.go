package gate

import (
	"github.com/Kamva/gutil"
	"github.com/Kamva/hexa"
)

// DefaultMiddlewareOptions is the default middleware's options.
type DefaultMiddlewareOptions struct {
	DenyGuest           bool
	DenyDeactivatedUser bool
	AllowRoot           bool

	RootPermName string
}

func DefaultMiddleware(o DefaultMiddlewareOptions) hexa.GateMiddleware {
	if o.RootPermName == "" {
		o.RootPermName = "root"
	}
	return func(policy hexa.GatePolicy) hexa.GatePolicy {
		return func(ctx hexa.Context, user hexa.User, resource interface{}) (b bool, err error) {
			permList := user.PermissionsList()
			// Check guest user
			if o.DenyGuest && user.Type() == hexa.UserTypeGuest {
				return false, nil
			}
			// Check user activation
			if o.DenyDeactivatedUser && !user.IsActive() { // TODO: How we should share same values?

				return false, nil
			}
			// check root
			if o.AllowRoot && gutil.Contains(permList, "root") { // TODO: replace root word with shared defined variable.
				return true, nil
			}
			return policy(ctx, user, resource)
		}
	}
}

// Assertion
var _ hexa.GateMiddleware = DefaultMiddleware(DefaultMiddlewareOptions{})
