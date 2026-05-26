// hostservice_capability_codec.go exposes organization and tenant capability host service payload codecs.
// Capability authorization remains outside this facade; declarations here are wire-level aliases only.

package protocol

import "lina-core/pkg/plugin/pluginbridge/internal/hostservice"

var (
	MarshalHostServiceCapabilityJSONResponse              = hostservice.MarshalHostServiceCapabilityJSONResponse
	UnmarshalHostServiceCapabilityJSONResponse            = hostservice.UnmarshalHostServiceCapabilityJSONResponse
	MarshalHostServiceCapabilityTenantRequest             = hostservice.MarshalHostServiceCapabilityTenantRequest
	UnmarshalHostServiceCapabilityTenantRequest           = hostservice.UnmarshalHostServiceCapabilityTenantRequest
	MarshalHostServiceCapabilityUserRequest               = hostservice.MarshalHostServiceCapabilityUserRequest
	UnmarshalHostServiceCapabilityUserRequest             = hostservice.UnmarshalHostServiceCapabilityUserRequest
	MarshalHostServiceCapabilityUsersRequest              = hostservice.MarshalHostServiceCapabilityUsersRequest
	UnmarshalHostServiceCapabilityUsersRequest            = hostservice.UnmarshalHostServiceCapabilityUsersRequest
	MarshalHostServiceCapabilityUserTenantRequest         = hostservice.MarshalHostServiceCapabilityUserTenantRequest
	UnmarshalHostServiceCapabilityUserTenantRequest       = hostservice.UnmarshalHostServiceCapabilityUserTenantRequest
	MarshalHostServiceCapabilityUserTenantSwitchRequest   = hostservice.MarshalHostServiceCapabilityUserTenantSwitchRequest
	UnmarshalHostServiceCapabilityUserTenantSwitchRequest = hostservice.UnmarshalHostServiceCapabilityUserTenantSwitchRequest
)
