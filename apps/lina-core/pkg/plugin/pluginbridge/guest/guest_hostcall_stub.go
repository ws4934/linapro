//go:build !wasip1

// This file provides the non-WASI raw host-call transport stub used by
// higher-level guest capability clients.

package guest

// InvokeHostService reports that generic guest host calls are unavailable.
func InvokeHostService(_ string, _ string, _ string, _ string, _ []byte) ([]byte, error) {
	return nil, ErrHostCallsUnavailable
}
