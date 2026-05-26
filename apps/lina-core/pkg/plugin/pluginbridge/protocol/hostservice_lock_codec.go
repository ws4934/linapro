// hostservice_lock_codec.go exposes distributed lock host service payload codecs.
// Lock ownership and renewal behavior remain implemented by host services; this facade only re-exports wire helpers.

package protocol

import "lina-core/pkg/plugin/pluginbridge/internal/hostservice"

var (
	MarshalHostServiceLockAcquireRequest    = hostservice.MarshalHostServiceLockAcquireRequest
	UnmarshalHostServiceLockAcquireRequest  = hostservice.UnmarshalHostServiceLockAcquireRequest
	MarshalHostServiceLockAcquireResponse   = hostservice.MarshalHostServiceLockAcquireResponse
	UnmarshalHostServiceLockAcquireResponse = hostservice.UnmarshalHostServiceLockAcquireResponse
	MarshalHostServiceLockRenewRequest      = hostservice.MarshalHostServiceLockRenewRequest
	UnmarshalHostServiceLockRenewRequest    = hostservice.UnmarshalHostServiceLockRenewRequest
	MarshalHostServiceLockRenewResponse     = hostservice.MarshalHostServiceLockRenewResponse
	UnmarshalHostServiceLockRenewResponse   = hostservice.UnmarshalHostServiceLockRenewResponse
	MarshalHostServiceLockReleaseRequest    = hostservice.MarshalHostServiceLockReleaseRequest
	UnmarshalHostServiceLockReleaseRequest  = hostservice.UnmarshalHostServiceLockReleaseRequest
)
