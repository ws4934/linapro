// Package valuecopy contains small copy helpers used by pluginhost wrappers.
package valuecopy

// Map returns a shallow copy of the given payload value map.
func Map(values map[string]interface{}) map[string]interface{} {
	if len(values) == 0 {
		return map[string]interface{}{}
	}
	cloned := make(map[string]interface{}, len(values))
	for key, value := range values {
		cloned[key] = value
	}
	return cloned
}
