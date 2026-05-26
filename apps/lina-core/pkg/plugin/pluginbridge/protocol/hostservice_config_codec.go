// hostservice_config_codec.go exposes configuration host service payload codecs.
// This file only preserves protocol aliases and does not own configuration lookup behavior.

package protocol

import "lina-core/pkg/plugin/pluginbridge/internal/hostservice"

var (
	MarshalHostServiceConfigKeyRequest      = hostservice.MarshalHostServiceConfigKeyRequest
	UnmarshalHostServiceConfigKeyRequest    = hostservice.UnmarshalHostServiceConfigKeyRequest
	MarshalHostServiceConfigValueResponse   = hostservice.MarshalHostServiceConfigValueResponse
	UnmarshalHostServiceConfigValueResponse = hostservice.UnmarshalHostServiceConfigValueResponse
)
