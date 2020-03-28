package gate

import (
	"github.com/Kamva/gutil"
	"github.com/Kamva/hexa"
)

func DefaultMiddleware(denyGuest, denyDeactivatedUser, allowsRoot bool) hexa.GateMiddleware {
	return func(policy hexa.GatePolicy) hexa.GatePolicy {
		return func(ctx hexa.Context, user hexa.User, resource interface{}) (b bool, err error) {
			permList := user.GetPermissionsList()
			// Check guest user
			if denyGuest && user.IsGuest() {
				return false, nil
			}
			// Check user activation
			if denyDeactivatedUser && !gutil.Contains(permList, "activated_account") { // TODO: How we should share same values?

				return false, nil
			}
			// check root
			if allowsRoot && gutil.Contains(permList, "root") { // TODO: replace root word with shared defined variable.
				return true, nil
			}
			return policy(ctx, user, resource)
		}
	}
}

// Assertion
var _ hexa.GateMiddleware = DefaultMiddleware(false,false,false,)
