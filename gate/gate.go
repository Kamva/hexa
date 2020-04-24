// gate package implements the hexa.Gate interface.
package gate

import (
	"github.com/Kamva/gutil"
	"github.com/Kamva/hexa"
)

type gate struct {
	m hexa.GateMiddleware
	p hexa.GatePolicy
}

func (g *gate) WithPolicy(p hexa.GatePolicy) hexa.Gate {
	return NewWithOptions(g.m, p)
}

func (g *gate) AllowsRoot(ctx hexa.Context) (bool, error) {
	return g.AllowsWithOptions(ctx, hexa.GateAllowsOptions{
		GatePolicy:        g.p,
		UserPermission:    "",
		ManagerPermission: "",
		Resource:          nil,
	})
}

func (g *gate) Allows(ctx hexa.Context, perm string, resource interface{}) (bool, error) {
	managerPerm, userPerm := g.extractPerms(perm)
	return g.AllowsWithOptions(ctx, hexa.GateAllowsOptions{
		GatePolicy:        g.p,
		UserPermission:    userPerm,
		ManagerPermission: managerPerm,
		Resource:          resource,
	})
}

func (g *gate) AllowsResource(ctx hexa.Context, resource interface{}) (bool, error) {
	return g.Allows(ctx, "", resource)
}

func (g *gate) AllowsWithOptions(c hexa.Context, options hexa.GateAllowsOptions) (bool, error) {
	return g.m(g.handler)(c, options)
}

func (g *gate) handler(c hexa.Context, options hexa.GateAllowsOptions) (bool, error) {
	user := c.User()
	allowsManager := g.allowsManager(user, options.ManagerPermission)
	allowsUser := g.allowsUser(user, options.UserPermission)
	policyIsOk, err := options.GatePolicy(c, user, )

	return allowsManager || (allowsUser && policyIsOk), err
}

// allowsManager check if user has specified permission as managing permission.
func (g *gate) allowsManager(user hexa.User, perm string) bool {
	if perm == "" {
		return false
	}
	return gutil.Contains(user.PermissionsList(), perm)
}

// allowsUser check if user has specified permission.
func (g *gate) allowsUser(user hexa.User, perm string) bool {
	if perm == "" {
		return true
	}
	return gutil.Contains(user.PermissionsList(), perm)
}

// extractPerms generate manager permission and user permission from provided
// permission.
// If perm is empty, returns empty.
// If perm start with "mgr:", returns perm,""
// return "mgr:{perm}",perm
func (g *gate) extractPerms(perm string) (managerPerm, userPerm string) {
	if len(perm) == 0 {
		return "", ""
	}
	mgrPerm := ManagerPerm(perm)
	if mgrPerm == perm {
		return perm, ""
	}
	return mgrPerm, perm
}

// New returns new instance of the Gate.
func New() hexa.Gate {
	return &gate{m: DefaultMiddleware(DefaultMiddlewareOptions{
		DenyGuest:           true,
		DenyDeactivatedUser: true,
		AllowRoot:           true,
	}), p: DefaultPolicy}
}

// New returns new instance of the Gate.
func NewWithOptions(m hexa.GateMiddleware, p hexa.GatePolicy) hexa.Gate {
	return &gate{m: m, p: p}
}

var _ hexa.Gate = &gate{}
var _ hexa.GateHandler = (&gate{}).handler
