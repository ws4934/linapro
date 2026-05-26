// This file defines the source-plugin route binding snapshot captured by the
// host during route registration for later governance and OpenAPI projection.

package pluginhost

import "lina-core/pkg/plugin/pluginhost/internal/routebinding"

// routeMethodAll is the normalized wildcard method marker used by route capture.
const routeMethodAll = routebinding.MethodAll

// SourceRouteBinding stores one plugin-owned route binding captured during host
// route registration.
type SourceRouteBinding struct {
	// PluginID is the owning source-plugin identifier.
	PluginID string
	// Method is the resolved HTTP method, such as GET or POST.
	Method string
	// Path is the resolved public route path registered on the host server.
	Path string
	// Handler is the bound route handler or bound object method.
	Handler interface{}
	// Documentable reports whether the handler uses GoFrame standard DTO routing
	// and can therefore be projected into the host OpenAPI document.
	Documentable bool
}

// Key returns the normalized uniqueness key of the binding.
func (b SourceRouteBinding) Key() string {
	return normalizeRouteMethod(b.Method) + " " + normalizeRoutePattern(b.Path)
}

// CloneSourceRouteBindings returns one detached copy of the given bindings.
func CloneSourceRouteBindings(bindings []SourceRouteBinding) []SourceRouteBinding {
	if len(bindings) == 0 {
		return []SourceRouteBinding{}
	}
	items := make([]SourceRouteBinding, len(bindings))
	copy(items, bindings)
	return items
}

// captureRouteBindings captures one or more route bindings from either a
// function handler or a controller object.
func captureRouteBindings(
	pluginID string,
	prefix string,
	explicitPattern string,
	explicitMethod string,
	handler interface{},
) []SourceRouteBinding {
	return toSourceRouteBindings(routebinding.Capture(pluginID, prefix, explicitPattern, explicitMethod, handler))
}

// normalizeRouteMethod canonicalizes an HTTP method and falls back to ALL when
// the input is empty.
func normalizeRouteMethod(method string) string {
	return routebinding.NormalizeRouteMethod(method)
}

// normalizeRoutePattern canonicalizes one route path for stable capture keys
// and path composition.
func normalizeRoutePattern(pattern string) string {
	return routebinding.NormalizeRoutePattern(pattern)
}

// joinRoutePatterns joins a group prefix and child pattern into one normalized
// public route path.
func joinRoutePatterns(prefix string, pattern string) string {
	return routebinding.JoinRoutePatterns(prefix, pattern)
}

// toSourceRouteBindings converts internal routebinding snapshots to the
// published pluginhost route binding contract.
func toSourceRouteBindings(bindings []routebinding.Binding) []SourceRouteBinding {
	if len(bindings) == 0 {
		return nil
	}
	items := make([]SourceRouteBinding, 0, len(bindings))
	for _, binding := range bindings {
		items = append(items, SourceRouteBinding{
			PluginID:     binding.PluginID,
			Method:       binding.Method,
			Path:         binding.Path,
			Handler:      binding.Handler,
			Documentable: binding.Documentable,
		})
	}
	return items
}
