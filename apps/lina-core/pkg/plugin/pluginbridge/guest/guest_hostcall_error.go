// This file defines shared guest host-call transport errors.

package guest

import "github.com/gogf/gf/v2/errors/gerror"

// ErrHostCallsUnavailable reports that guest host calls are unavailable in
// non-WASI builds.
var ErrHostCallsUnavailable = gerror.New(
	"pluginbridge guest host-call transport is only available for wasip1 builds",
)
