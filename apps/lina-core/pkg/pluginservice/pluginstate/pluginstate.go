// Package pluginstate exposes source-plugin enablement lookups without
// publishing host-internal plugin service packages to plugin implementations.
package pluginstate

import "lina-core/pkg/pluginservice/contract"

// serviceAdapter bridges host plugin state into the published plugin contract.
type serviceAdapter struct {
	service contract.PluginStateService
}

// New creates and returns a plugin state service backed by the given service.
func New(service contract.PluginStateService) contract.PluginStateService {
	return &serviceAdapter{service: service}
}
