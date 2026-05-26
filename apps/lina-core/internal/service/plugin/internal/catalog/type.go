// This file defines the plugin type enumeration and normalization helpers.

package catalog

import "strings"

// PluginType defines the recognized plugin types.
type PluginType string

// Canonical plugin type values supported by the host.
const (
	// TypeSource identifies a compiled-in source plugin.
	TypeSource PluginType = "source"
	// TypeDynamic identifies a runtime-loaded WASM plugin.
	TypeDynamic PluginType = "dynamic"
)

// String returns the canonical type value.
func (t PluginType) String() string { return string(t) }

// NormalizeType converts a raw type string to the canonical PluginType.
func NormalizeType(value string) PluginType {
	switch strings.TrimSpace(strings.ToLower(value)) {
	case TypeSource.String():
		return TypeSource
	case TypeDynamic.String():
		return TypeDynamic
	default:
		return PluginType(strings.TrimSpace(strings.ToLower(value)))
	}
}

// IsSupportedType reports whether the given type string is a recognized plugin type.
func IsSupportedType(value string) bool {
	t := NormalizeType(value)
	return t == TypeSource || t == TypeDynamic
}
