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
	return func(h hexa.GateHandler) hexa.GateHandler {
		return func(ctx hexa.Context, handlerOptions hexa.GateAllowsOptions) (b bool, err error) {
			user := ctx.User()
			permList := user.PermissionsList()
			// Check guest user
			if o.DenyGuest && user.Type() == hexa.UserTypeGuest {
				return false, nil
			}
			// Check user activation
			if o.DenyDeactivatedUser && !user.IsActive() {
				return false, nil
			}
			// check root
			if o.AllowRoot && gutil.Contains(permList, "root") {
				return true, nil
			}
			return h(ctx, handlerOptions)
		}
	}
}
