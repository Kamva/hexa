package hexa

type (
	// GatePolicy is the ABAC policy
	GatePolicy func(ctx Context, resource interface{}) (bool, error)

	// GateHandler handles gate authorization requests.
	GateHandler func(Context, GateAllowsOptions) (bool, error)

	// GateMiddleware is the Gate middleware (this middleware will check some conditions before run policy).
	// e.g: if user has root permission, so it returns true and dont call to the ABAC policy.
	GateMiddleware func(GateHandler) GateHandler

	// GateAllowsOptions contains options which provide to the gate on authorization
	GateAllowsOptions struct {
		GatePolicy

		UserPermission    string
		ManagerPermission string
		Resource          interface{}
	}

	// Gate defines the gate interface.
	Gate interface {
		// WithPolicy returns new instance of gate with provided policy.
		WithPolicy(p GatePolicy) Gate

		// AllowsRoot returns true just when user has the "root" permission.
		// This function's behavior is relative to your Gate middleware, because
		// gate middleware should take care of the root permission.
		AllowsRoot(Context) (bool, error)

		// Allows returns true when either one of these situations is true:
		// - user has "mgr:{permission}" permission
		// - user has permission {perm} and policy function returns true.
		//
		// Describe as expression: (has_permission "mgr:{perm}" || has_perm {perm} && policy(resource))
		//
		// Notes:
		// - If user is guest, returns false(default).
		// - If user is not active(user's "IsActive" method returns false), returns false(default)
		// - If user is root, returns true (default).
		//
		//
		// Permission Name :
		// - manager(back-office) permissions has "mgr" prefix, for example if you pass "create_post" permission, we
		// check the "mgr:create_post" permission also and if user hvs mgr:create_post permission we return
		// true, otherwise check user policy is true and user has create_post permission.
		//
		// Examples:
		// - To permit to the manager:
		//   Allows("mgr:review_all_payments",nil) // return true if manager has mgr:review_all_payments permission.
		//
		// - To permit to just manager with permission "mgr:create_post" or (user with permission "create_post" and
		//   also check provided user_id in the payload is the current user's id), call:
		// 	 gate.Allows("edit_post",gate.NewEmptyResourceWithOwner(userID))
		//
		// - To permit to just manager with permission "mgr:edit_post" or user with permission "edit_post" call:
		// 	 gate.WithPolicy(TruePolicy).Allows("edit_post",nil)
		//
		// - To permit to just manager with permission "mgr:edit_post" or (user with permission "edit_post"
		//   and permit by policy) call:  gate.WithPolicy(YourPolicy).Allows("edit_post",post)
		//
		// - To permit user with specific resource(e.g product) without checking any permission (just checking by policy function)
		//  call: gate.fromPolicy(YourPolicy).Allows("",product)
		//  or call: gate.WithPolicy(YourPolicy).AllowsResource(product)
		Allows(ctx Context, perm string, resource interface{}) (bool, error)

		// AllowsResource just checks policy.
		// This is alias for call to `gate.Allows(ctx,"",resource)`
		AllowsResource(ctx Context, resource interface{}) (bool, error)

		// Check user Allows to do something. Get full options.
		AllowsWithOptions(Context, GateAllowsOptions) (bool, error)
	}
)
