package gate

import (
	"fmt"
	"strings"
)

// managerPermPrefix is the permission prefix for
// manager permissions.
const managerPermPrefix = "mgr"

// ManagerPerm returns management permission type of a permission.
// e.g `update_post` permission let post owner to update update his
// post, if a manager has "mgr:update_post", he/she can also update
// that post.
func ManagerPerm(perm string) string {
	if strings.HasPrefix(perm, "mgr:") {
		return perm
	}
	return fmt.Sprintf("%s:%s", managerPermPrefix, perm)
}
