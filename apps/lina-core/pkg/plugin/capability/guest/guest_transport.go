// This file binds guest capability clients to the raw pluginbridge guest
// transport without re-exporting bridge protocol DTOs.

package guest

import (
	"github.com/gogf/gf/v2/errors/gerror"

	bridgeguest "lina-core/pkg/plugin/pluginbridge/guest"
)

var (
	// ErrHostCallsUnavailable reports that capability guest host-service
	// clients are unavailable in non-WASI builds.
	ErrHostCallsUnavailable = gerror.New(
		"capability guest host-service clients are only available for wasip1 builds",
	)
)

// invokeHostService dispatches one structured host-service request through the
// raw pluginbridge guest transport.
func invokeHostService(service string, method string, resourceRef string, table string, payload []byte) ([]byte, error) {
	response, err := bridgeguest.InvokeHostService(service, method, resourceRef, table, payload)
	if gerror.Is(err, bridgeguest.ErrHostCallsUnavailable) {
		return nil, ErrHostCallsUnavailable
	}
	return response, err
}
