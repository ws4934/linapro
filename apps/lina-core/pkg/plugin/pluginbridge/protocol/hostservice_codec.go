// hostservice_codec.go exposes generic host service envelope and value codecs.
// Domain-specific request and response codecs are split into neighboring files to keep each protocol surface small.

package protocol

import "lina-core/pkg/plugin/pluginbridge/internal/hostservice"

var (
	MarshalHostServiceRequestEnvelope   = hostservice.MarshalHostServiceRequestEnvelope
	UnmarshalHostServiceRequestEnvelope = hostservice.UnmarshalHostServiceRequestEnvelope
	MarshalHostServiceValueResponse     = hostservice.MarshalHostServiceValueResponse
	UnmarshalHostServiceValueResponse   = hostservice.UnmarshalHostServiceValueResponse
)
