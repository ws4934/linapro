// Package hostconfig exposes a small, whitelisted public host config reader for
// plugins. It is intentionally separate from plugin config so plugins cannot
// scan or read the host's complete static configuration tree.
package hostconfig

import (
	hostconfigsvc "lina-core/internal/service/config"
	"lina-core/pkg/plugin/capability/contract"
)

// serviceAdapter reads public host config keys from the host config service.
type serviceAdapter struct {
	configSvc hostconfigsvc.Service
}

// New creates a public host config reader backed by the host config service.
func New(configSvc hostconfigsvc.Service) contract.HostConfigService {
	return &serviceAdapter{configSvc: configSvc}
}
