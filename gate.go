package hexa

type (
	// GatePolicy is the ABAC policy
	GatePolicy func(ctx Context, user User, resource interface{}) (bool, error)

	// GateMiddleware is the Gate middleware.
	GateMiddleware func(GatePolicy) GatePolicy

	GateAllowsOptions struct {
		GateMiddleware // Gate can has just one GateMiddleware middleware.
		GatePolicy

		UserPermission    string
		ManagerPermission string
		Resource          interface{}
	}

	// Gate defines the gate interface.
	Gate interface {
		// FromPolicy returns new instance of gate with provided policy.
		FromPolicy(p GatePolicy) Gate

		// Allows returns true either when one of this situations is true:
		// - user has "mgr:{permission}" permission
		// - user has permission {perm} and policy fn returns true.
		//
		// Match expression: (has_permission "mgr:{perm}" || has_perm {perm} && policy(resource))
		//
		// Notes:
		// - If user is guest, returns false(default).
		// - If user is not active(dont have "active_user" permission), returns false(default)
		// - If user is root, returns true (default).
		//
		//
		// Permission Name :
		// - manager(back-office) permissions has "mgr" prefix, for example if you pass "create_post" permission, we
		// check the "mgr:create_post" permission also and if user hvs mgr:create_post permission we return
		// true, otherwise check user policy is true and user has create_post permission.
		//
		// Some tricks:
		// - To permit to just the manager, do call like this:
		//   Allows("mgr:review_all_payments") // return true if manager has mgr:review_all_payments permission.
		//
		// - To permit to just manager with permission "mgr:edit_post" or user with permission "edit_post" call:
		// 	 gate.FromPolicy(DisabledPolicy).Allows("edit_post")
		//
		// - To permit to just manager with permission "mgr:edit_post" or (user with permission "edit_post"
		//   and permit by policy) call:  Allows("edit_post",post)
		Allows(ctx Context, perm string, resource interface{}) (bool, error)

		// AllowsResource just checks policy.
		AllowsResource(ctx Context, resource interface{}) (bool, error)

		// Check user Allows to do something. Get full options.
		AllowsWithOptions(Context, GateAllowsOptions) (bool, error)
	}
)
