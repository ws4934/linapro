// Package hostconfig exposes read-only host configuration access for source
// plugins. It is intentionally separate from plugin config because plugin
// config is scoped to one plugin's own runtime configuration files.
package hostconfig

import (
	"context"

	"github.com/gogf/gf/v2/container/gvar"

	hostconfigsvc "lina-core/internal/service/config"
	"lina-core/pkg/plugin/capability/contract"
)

// rawConfigReader is implemented by the host-owned config service. It keeps
// this adapter dependent on the startup-injected config instance instead of
// reaching around the service graph for global configuration.
type rawConfigReader interface {
	GetRaw(ctx context.Context, key string) (*gvar.Var, error)
}

// serviceAdapter reads individual host config keys from the host config service.
type serviceAdapter struct {
	configSvc hostconfigsvc.Service
}

// New creates a host config reader backed by the host config service.
func New(configSvc hostconfigsvc.Service) contract.HostConfigService {
	return &serviceAdapter{configSvc: configSvc}
}
