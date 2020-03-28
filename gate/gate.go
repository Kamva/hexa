// gate package implements the hexa.Gate interface.
package gate

import (
	"fmt"
	"github.com/Kamva/gutil"
	"github.com/Kamva/hexa"
	"strings"
)

type gate struct {
	m hexa.GateMiddleware
	p hexa.GatePolicy
}

func (g *gate) FromPolicy(p hexa.GatePolicy) hexa.Gate {
	return NewWithOptions(g.m, p)
}

func (g *gate) Allows(ctx hexa.Context, perm string, resource interface{}) (bool, error) {
	managerPerm, userPerm := g.extractPerms(perm)
	return g.AllowsWithOptions(ctx, hexa.GateAllowsOptions{
		GateMiddleware:    g.m,
		GatePolicy:        g.p,
		UserPermission:    userPerm,
		ManagerPermission: managerPerm,
		Resource:          resource,
	})
}

func (g *gate) AllowsWithOptions(c hexa.Context, options hexa.GateAllowsOptions) (bool, error) {
	user := c.User()
	allowsManager := g.allowsManager(user, options.ManagerPermission)
	allowsUser := g.allowsUser(user, options.UserPermission)
	policyIsOk, err := options.GateMiddleware(options.GatePolicy)(c, user, options.Resource)

	return allowsManager || (allowsUser && policyIsOk), err
}

// allowsManager check if user has specified permission as managing permission.
func (g *gate) allowsManager(user hexa.User, perm string) bool {
	if perm == "" {
		return false
	}
	return gutil.Contains(user.GetPermissionsList(), perm)
}

// allowsUser check if user has specified permission.
func (g *gate) allowsUser(user hexa.User, perm string) bool {
	if perm == "" {
		return true
	}
	return gutil.Contains(user.GetPermissionsList(), perm)
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
	if strings.HasPrefix(perm, "mgr:") {
		return perm, ""
	}
	return fmt.Sprintf("mgr:%s", perm), perm
}

// New returns new instance of the Gate.
func New() hexa.Gate {
	return &gate{m: DefaultMiddleware(true, true, true), p: DefaultPolicy}
}

// New returns new instance of the Gate.
func NewWithOptions(m hexa.GateMiddleware, p hexa.GatePolicy) hexa.Gate {
	return &gate{m: m, p: p}
}

var _ hexa.Gate = &gate{}
